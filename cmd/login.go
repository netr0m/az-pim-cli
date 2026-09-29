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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in for the Entra roles/groups commands, replacing any cached account",
	Run: func(cmd *cobra.Command, args []string) {
		requireGraphCredentials()
		if err := pim.ClearGraphTokenCache(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		token := pim.GetGraphAccessToken([]string{pim.GRAPH_DEFAULT_SCOPE}, AzureClientInstance)
		userInfo := pim.GetUserInfo(token)
		slog.Info("Signed in", "email", userInfo.Email)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
