/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"log"

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
		if Token == "" {
			log.Fatalf("Listing eligible groups requires providing a token manually due to restrictions in token permissions. Consult the docs for more information.")
		}
		subjectId := pim.GetUserInfo(Token).ObjectId
		eligibleGroupAssignments := pim.GetEligibleGroupAssignments(Token, subjectId)
		utils.PrintEligibleGroups(eligibleGroupAssignments)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listGroupCmd)
}
