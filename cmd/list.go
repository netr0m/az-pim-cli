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
		token := pim.GetPIMAccessTokenAzureCLI(pim.AZ_PIM_SCOPE)

		eligibleRoleAssignments := pim.GetEligibleRoleAssignments(token)
		utils.PrintEligibleRoles(eligibleRoleAssignments)
	},
}

var listGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g", "grp", "groups"},
	Short:   "Query Azure PIM for eligible group assignments",
	Run: func(cmd *cobra.Command, args []string) {
		subjectId := pim.GetUserInfo(pimGroupsToken).ObjectId
		eligibleGroupAssignments := pim.GetEligibleGroupAssignments(pimGroupsToken, subjectId)
		utils.PrintEligibleGroups(eligibleGroupAssignments)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listGroupCmd)

	listGroupCmd.PersistentFlags().StringVarP(&pimGroupsToken, "token", "t", "", "An access token for the PIM Groups API (required). Consult the README for more information.")
	listGroupCmd.MarkPersistentFlagRequired("token") //nolint:errcheck
}
