/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"github.com/netr0m/az-pim-cli/pkg/pim"
	"github.com/netr0m/az-pim-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "Query Azure PIM for eligible role assignments",
	Run: func(cmd *cobra.Command, args []string) {
		token := pim.GetPIMAccessTokenAzureCLI()

		eligibleRoleAssignments := pim.GetEligibleRoleAssignments(token)
		utils.PrintEligibleRoles(eligibleRoleAssignments)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
