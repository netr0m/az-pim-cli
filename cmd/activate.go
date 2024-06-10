/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/netr0m/az-pim-cli/pkg/pim"
	"github.com/netr0m/az-pim-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var subscriptionName string
var subscriptionPrefix string
var roleName string
var duration int
var reason string

var activateCmd = &cobra.Command{
	Use:     "activate",
	Aliases: []string{"a", "ac", "act"},
	Short:   "Sends a request to Azure PIM to activate the given role",
	Run: func(cmd *cobra.Command, args []string) {
		if subscriptionName == "" && subscriptionPrefix == "" {
			log.Fatalf("Missing required parameter: You must specify either 'subscription-name' or 'subscription-prefix'.")
		}
		token := pim.GetPIMAccessTokenAzureCLI(pim.AZ_PIM_SCOPE)
		subjectId := pim.GetUserInfo(token).ObjectId

		eligibleRoleAssignments := pim.GetEligibleRoleAssignments(token)
		roleAssignment := utils.GetRoleAssignment(subscriptionName, subscriptionPrefix, roleName, eligibleRoleAssignments)

		log.Printf(
			"Activating role '%s' in subscription '%s' with reason '%s'",
			roleAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName,
			roleAssignment.Properties.ExpandedProperties.Scope.DisplayName,
			reason,
		)

		requestResponse := pim.RequestRoleAssignment(subjectId, roleAssignment, duration, reason, token)
		log.Printf("The role '%s' in '%s' is now %s", roleAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName, roleAssignment.Properties.ExpandedProperties.Scope.DisplayName, requestResponse.Properties.Status)
	},
}

var activateGroupCmd = &cobra.Command{
	Use:     "group [group name]",
	Aliases: []string{"g", "grp", "groups"},
	Short:   "Sends a request to Azure PIM to activate the given group",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		if len(args) == 0 {
			comps = cobra.AppendActiveHelp(comps, "Please specify the name of the group")
		} else if len(args) == 1 {
			comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
		} else {
			comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
		}
		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("activate group needs a name for the group"))
		}
		if Token == "" {
			log.Fatalf("Activating a group requires providing a token manually due to restrictions in token permissions. Consult the docs for more information.")
		}
		subjectId := pim.GetUserInfo(Token).ObjectId

		eligibleGroupAssignments := pim.GetEligibleGroupAssignments(Token, subjectId)
		groupAssignment := utils.GetGroupAssignment(args[0], roleName, eligibleGroupAssignments)

		log.Printf(
			"Activating role '%s' for group '%s' with reason '%s'",
			groupAssignment.RoleDefinition.DisplayName,
			groupAssignment.RoleDefinition.Resource.DisplayName,
			reason,
		)

		requestResponse := pim.RequestGroupAssignment(subjectId, groupAssignment, duration, reason, Token)
		log.Printf("The role '%s' for group '%s' is now %s", groupAssignment.RoleDefinition.DisplayName, groupAssignment.RoleDefinition.Resource.DisplayName, requestResponse.AssignmentState)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.AddCommand(activateGroupCmd)

	// Flags
	activateCmd.PersistentFlags().StringVarP(&subscriptionName, "subscription-name", "s", "", "The name of the subscription to activate")
	activateCmd.PersistentFlags().StringVarP(&subscriptionPrefix, "subscription-prefix", "p", "", "The name prefix of the subscription to activate (e.g. 'S399'). Alternative to 'subscription-name'.")
	activateCmd.PersistentFlags().StringVarP(&roleName, "role-name", "r", "", "Specify the role to activate, if multiple roles are found for a subscription (e.g. 'Owner' and 'Contributor')")
	activateCmd.PersistentFlags().IntVarP(&duration, "duration", "d", pim.DEFAULT_DURATION_MINUTES, "Duration in minutes that the role should be activated for")
	activateCmd.PersistentFlags().StringVar(&reason, "reason", pim.DEFAULT_REASON, "Reason for the activation")
}
