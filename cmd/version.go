/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	goVersion = "unknown"

	// populated by goreleaser during build
	version = "unknown"
	commit  = "?"
	date    = ""
)

var includeBuildInfo bool

type BuildInfo struct {
	GoVersion string `json:"goVersion"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
}

func (b BuildInfo) String() string {
	return fmt.Sprintf("az-pim-cli version %s (built with %s from %s on %s)",
		b.Version, b.GoVersion, b.Commit, b.Date)
}

func printVersion(w io.Writer, info BuildInfo) {
	fmt.Fprintln(w, info.String())
}

func createBuildInfo() BuildInfo {
	info := BuildInfo{
		Commit:    commit,
		Version:   version,
		GoVersion: goVersion,
		Date:      date,
	}

	buildInfo, available := debug.ReadBuildInfo()
	if !available {
		return info
	}

	info.GoVersion = buildInfo.GoVersion

	if date != "" {
		return info
	}

	info.Version = buildInfo.Main.Version

	var revision = "unknown"
	var modified = "?"
	for _, setting := range buildInfo.Settings {
		// vcs subkeys are available during build
		switch setting.Key {
		case "vcs.time":
			info.Date = setting.Value
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}

	if info.Date == "" {
		info.Date = "(unknown)"
	}

	info.Commit = fmt.Sprintf("(%s, modified: %s, mod sum: %q)", revision, modified, buildInfo.Main.Sum)

	return info
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of az-pim-cli",
	Run: func(cmd *cobra.Command, args []string) {
		if includeBuildInfo {
			debugInfo, ok := debug.ReadBuildInfo()
			if !ok {
				fmt.Fprintln(os.Stderr, "Failed to read build info")
			}
			fmt.Fprintln(os.Stdout, debugInfo)
		}
		info := createBuildInfo()
		printVersion(os.Stdout, info)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.PersistentFlags().BoolVarP(&includeBuildInfo, "debug", "d", false, "Include build information")
}
