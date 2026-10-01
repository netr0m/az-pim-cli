/*
Copyright © 2026 netr0m <netr0m@pm.me>
*/
package pim

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type fileCache struct {
	path string
}

func (f *fileCache) Replace(_ context.Context, unmarshaler cache.Unmarshaler, _ cache.ReplaceHints) error {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read token cache: %w", err)
	}

	if err := unmarshaler.Unmarshal(data); err != nil {
		return fmt.Errorf("failed to unmarshal token cache: %w", err)
	}
	return nil
}

func (f *fileCache) Export(_ context.Context, marshaler cache.Marshaler, _ cache.ExportHints) error {
	data, err := marshaler.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal token cache: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return fmt.Errorf("failed to create token cache directory: %w", err)
	}
	if err := os.WriteFile(f.path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write token cache: %w", err)
	}
	return nil
}

func graphTokenCachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user cache directory: %w", err)
	}
	return filepath.Join(dir, "az-pim-cli", "msal_cache.json"), nil
}

func ClearGraphTokenCache() error {
	path, err := graphTokenCachePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove token cache: %w", err)
	}
	return nil
}

func acquireGraphToken(clientID, tenantID string, scopes []string, cachePath, authorityHost string) (string, error) {
	client, err := public.New(
		clientID,
		public.WithAuthority(fmt.Sprintf("https://%s/%s", authorityHost, tenantID)),
		public.WithCache(&fileCache{path: cachePath}),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create MSAL client: %w", err)
	}

	ctx := context.Background()

	if accounts, err := client.Accounts(ctx); err == nil && len(accounts) > 0 {
		result, err := client.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(accounts[0]))
		if err == nil {
			return result.AccessToken, nil
		}
	}

	deviceCode, err := client.AcquireTokenByDeviceCode(ctx, scopes)
	if err != nil {
		return "", fmt.Errorf("failed to start device code login: %w", err)
	}

	fmt.Println(deviceCode.Result.Message) //nolint:forbidigo

	result, err := deviceCode.AuthenticationResult(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to complete device code login: %w", err)
	}

	return result.AccessToken, nil
}
