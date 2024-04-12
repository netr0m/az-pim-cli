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
		token := pim.GetPIMAccessTokenAzureCLI()
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

func init() {
	rootCmd.AddCommand(activateCmd)

	// Flags
	activateCmd.PersistentFlags().StringVarP(&subscriptionName, "subscription-name", "s", "", "The name of the subscription to activate")
	activateCmd.PersistentFlags().StringVarP(&subscriptionPrefix, "subscription-prefix", "p", "", "The name prefix of the subscription to activate (e.g. 'S399'). Alternative to 'subscription-name'.")
	activateCmd.PersistentFlags().StringVarP(&roleName, "role-name", "r", "", "Specify the role to activate, if multiple roles are found for a subscription (e.g. 'Owner' and 'Contributor')")
	activateCmd.PersistentFlags().IntVarP(&duration, "duration", "d", pim.DEFAULT_DURATION_MINUTES, "Duration in minutes that the role should be activated for")
	activateCmd.PersistentFlags().StringVar(&reason, "reason", pim.DEFAULT_REASON, "Reason for the activation")
}
