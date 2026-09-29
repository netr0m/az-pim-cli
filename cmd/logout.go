/*
Copyright © 2026 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/netr0m/az-pim-cli/pkg/pim"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear the cached sign-in used for the Entra roles/groups commands",
	Run: func(cmd *cobra.Command, args []string) {
		if err := pim.ClearGraphTokenCache(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("Signed out")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
