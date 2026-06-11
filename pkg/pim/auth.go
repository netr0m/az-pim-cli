/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/netr0m/az-pim-cli/pkg/common"
	"github.com/pkg/browser"
)

// Name of the token cache file (stored in the user's home directory)
const TOKEN_CACHE_FILE_NAME string = ".az-pim-cli.cache.json"

// tokenCachePath returns the path used to persist the MSAL token cache.
func tokenCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fall back to the current working directory if the home dir is unavailable
		return TOKEN_CACHE_FILE_NAME
	}
	return filepath.Join(home, TOKEN_CACHE_FILE_NAME)
}

// fileTokenCache is a pure-Go implementation of cache.ExportReplace that
// persists the MSAL token cache (including the refresh token) to a file with
// 0600 permissions. This avoids the cgo dependency of the OS keychain and works
// identically across platforms, including headless/SSH/container environments.
type fileTokenCache struct {
	path string
	mu   sync.Mutex
}

func (c *fileTokenCache) Replace(ctx context.Context, cacheData cache.Unmarshaler, hints cache.ReplaceHints) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return cacheData.Unmarshal(data)
}

func (c *fileTokenCache) Export(ctx context.Context, cacheData cache.Marshaler, hints cache.ExportHints) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := cacheData.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}

// getGraphToken acquires a Microsoft Graph token for the configured app
// registration using the device code flow, reusing a cached token (with silent
// refresh) when available so the user is only prompted occasionally.
func (c AzureClient) getGraphToken(scope string) string {
	if c.ClientID == "" {
		_error := common.Error{
			Operation: "getGraphToken",
			Message:   "no client ID configured. PIM for Groups and Entra roles requires a custom app registration; set --client-id and --tenant-id (or PIM_CLIENTID/PIM_TENANTID)",
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	accessor := &fileTokenCache{path: tokenCachePath()}
	app, err := public.New(c.ClientID, public.WithAuthority(c.Authority), public.WithCache(accessor))
	if err != nil {
		_error := common.Error{Operation: "getGraphToken", Message: err.Error(), Err: err}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	scopes := []string{scope}

	// Try to use a cached token (or silently refresh) for an existing account
	if accounts, aerr := app.Accounts(ctx); aerr == nil {
		for _, account := range accounts {
			if c.TenantID != "" && account.Realm != "" && account.Realm != c.TenantID {
				continue
			}
			if result, serr := app.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(account)); serr == nil {
				return result.AccessToken
			}
		}
	}

	// Acquire a new token via an interactive sign-in. The browser-based
	// (authorization code + PKCE) flow is used by default, since the device
	// code flow is commonly blocked by Conditional Access policies. The device
	// code flow remains available via --device-code for headless environments.
	var result public.AuthResult
	if c.UseDeviceCode {
		result, err = acquireTokenByDeviceCode(ctx, app, scopes)
	} else {
		result, err = acquireTokenInteractive(ctx, app, scopes)
	}
	if err != nil {
		_error := common.Error{Operation: "getGraphToken", Message: err.Error(), Status: "401", Err: err}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	return result.AccessToken
}

// acquireTokenInteractive opens the system browser and completes an
// authorization code (PKCE) sign-in via a loopback redirect. The browser is
// opened via openURL, which honors the BROWSER environment variable.
func acquireTokenInteractive(ctx context.Context, app public.Client, scopes []string) (public.AuthResult, error) {
	_, _ = fmt.Fprintln(os.Stderr, "Opening a browser to sign in...")
	return app.AcquireTokenInteractive(ctx, scopes, public.WithOpenURL(openURL))
}

// openURL opens the given URL in a browser. If the BROWSER environment variable
// is set, it is honored (following the common convention also used by Python's
// webbrowser module); otherwise the platform default browser is used.
func openURL(url string) error {
	browserEnv := os.Getenv("BROWSER")
	if browserEnv == "" {
		return browser.OpenURL(url)
	}

	var lastErr error
	for _, entry := range strings.Split(browserEnv, string(os.PathListSeparator)) {
		name, args, ok := browserCommand(entry, url)
		if !ok {
			continue
		}
		cmd := exec.Command(name, args...) // #nosec G204 -- command is taken from the user's BROWSER environment variable by design
		if err := cmd.Start(); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	// BROWSER was set but contained no usable entries; fall back to the default.
	return browser.OpenURL(url)
}

// browserCommand parses a single BROWSER entry into a command name and its
// arguments. "%s" in the entry is replaced with the URL; if no "%s" is present,
// the URL is appended as the final argument. A trailing "&" (background marker)
// is ignored, since the browser is always started without waiting. The returned
// bool is false when the entry is empty.
func browserCommand(entry string, url string) (string, []string, bool) {
	entry = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(entry), "&"))
	fields := strings.Fields(entry)
	if len(fields) == 0 {
		return "", nil, false
	}

	substituted := false
	args := make([]string, 0, len(fields))
	for _, field := range fields[1:] {
		if strings.Contains(field, "%s") {
			field = strings.ReplaceAll(field, "%s", url)
			substituted = true
		}
		args = append(args, field)
	}
	if !substituted {
		args = append(args, url)
	}

	return fields[0], args, true
}

// acquireTokenByDeviceCode prints a code + URL for the user to enter in a
// browser on any device. Useful for headless/SSH environments, but may be
// blocked by Conditional Access policies.
func acquireTokenByDeviceCode(ctx context.Context, app public.Client, scopes []string) (public.AuthResult, error) {
	deviceCode, err := app.AcquireTokenByDeviceCode(ctx, scopes)
	if err != nil {
		return public.AuthResult{}, err
	}
	// Print the sign-in instructions (code + URL) to stderr
	_, _ = fmt.Fprintln(os.Stderr, deviceCode.Result.Message)
	return deviceCode.AuthenticationResult(ctx)
}
