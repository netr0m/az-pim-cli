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
	Run:     func(cmd *cobra.Command, args []string) {},
}

var listResourceCmd = &cobra.Command{
	Use:     "resource",
	Aliases: []string{"r", "res", "resource", "resources", "sub", "subs", "subscriptions"},
	Short:   "Query Azure PIM for eligible resource assignments (azure resources)",
	Run: func(cmd *cobra.Command, args []string) {
		token := pim.GetPIMAccessTokenAzureCLI(pim.AZ_PIM_SCOPE)

		eligibleResourceAssignments := pim.GetEligibleResourceAssignments(token)
		utils.PrintEligibleResources(eligibleResourceAssignments)
	},
}

var listGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g", "grp", "groups"},
	Short:   "Query Azure PIM for eligible group assignments",
	Run: func(cmd *cobra.Command, args []string) {
		subjectId := pim.GetUserInfo(pimGovernanceRoleToken).ObjectId
		eligibleGroupAssignments := pim.GetEligibleGovernanceRoleAssignments(pim.ROLE_TYPE_AAD_GROUPS, subjectId, pimGovernanceRoleToken)
		utils.PrintEligibleGovernanceRoles(eligibleGroupAssignments)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listResourceCmd)
	listCmd.AddCommand(listGroupCmd)

	listGroupCmd.PersistentFlags().StringVarP(&pimGovernanceRoleToken, "token", "t", "", "An access token for the PIM 'Entra Roles' and 'Groups' API (required). Consult the README for more information.")
	listGroupCmd.MarkPersistentFlagRequired("token") //nolint:errcheck
}
