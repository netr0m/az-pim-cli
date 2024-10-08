/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"log"
	"os"

	"github.com/netr0m/az-pim-cli/pkg/pim"
	"github.com/netr0m/az-pim-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var name string
var prefix string
var roleName string
var duration int
var reason string
var ticketSystem string
var ticketNumber string
var dryRun bool

var activateCmd = &cobra.Command{
	Use:     "activate",
	Aliases: []string{"a", "ac", "act"},
	Short:   "Send a request to Azure PIM to activate a role assignment",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var activateResourceCmd = &cobra.Command{
	Use:     "resource",
	Aliases: []string{"r", "res", "resource", "resources", "sub", "subs", "subscriptions"},
	Short:   "Sends a request to Azure PIM to activate the given resource (azure resources)",
	Run: func(cmd *cobra.Command, args []string) {
		token := pim.GetPIMAccessTokenAzureCLI(pim.AZ_PIM_SCOPE)
		subjectId := pim.GetUserInfo(token).ObjectId

		eligibleResourceAssignments := pim.GetEligibleResourceAssignments(token)
		resourceAssignment := utils.GetResourceAssignment(name, prefix, roleName, eligibleResourceAssignments)

		log.Printf(
			"Activating role '%s' for resource '%s' with reason '%s' (ticket: %s [%s])",
			resourceAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName,
			resourceAssignment.Properties.ExpandedProperties.Scope.DisplayName,
			reason,
			ticketNumber,
			ticketSystem,
		)

		if dryRun {
			log.Printf("Skipping activation due to 'dry-run'.")
			os.Exit(0)
		}
		requestResponse := pim.RequestResourceAssignment(subjectId, resourceAssignment, duration, reason, ticketSystem, ticketNumber, token)
		log.Printf("The role '%s' in '%s' is now %s", resourceAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName, resourceAssignment.Properties.ExpandedProperties.Scope.DisplayName, requestResponse.Properties.Status)
	},
}

var activateGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g", "grp", "groups"},
	Short:   "Sends a request to Azure PIM to activate the given group",
	Run: func(cmd *cobra.Command, args []string) {
		subjectId := pim.GetUserInfo(pimGovernanceRoleToken).ObjectId

		eligibleGroupAssignments := pim.GetEligibleGovernanceRoleAssignments(pim.ROLE_TYPE_AAD_GROUPS, subjectId, pimGovernanceRoleToken)
		groupAssignment := utils.GetGovernanceRoleAssignment(name, prefix, roleName, eligibleGroupAssignments)

		log.Printf(
			"Activating role '%s' for group '%s' with reason '%s' (ticket: %s [%s])",
			groupAssignment.RoleDefinition.DisplayName,
			groupAssignment.RoleDefinition.Resource.DisplayName,
			reason,
			ticketNumber,
			ticketSystem,
		)

		if dryRun {
			log.Printf("Skipping activation due to 'dry-run'.")
			os.Exit(0)
		}
		requestResponse := pim.RequestGovernanceRoleAssignment(subjectId, pim.ROLE_TYPE_AAD_GROUPS, groupAssignment, duration, reason, ticketSystem, ticketNumber, pimGovernanceRoleToken)
		log.Printf("The role '%s' for group '%s' is now %s", groupAssignment.RoleDefinition.DisplayName, groupAssignment.RoleDefinition.Resource.DisplayName, requestResponse.AssignmentState)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.AddCommand(activateResourceCmd)
	activateCmd.AddCommand(activateGroupCmd)

	// Flags
	activateCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "The name of the resource to activate")
	activateCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "The name prefix of the resource to activate (e.g. 'S399'). Alternative to 'name'.")
	activateCmd.PersistentFlags().StringVarP(&roleName, "role", "r", "", "Specify the role to activate, if multiple roles are found for a resource (e.g. 'Owner' and 'Contributor')")
	activateCmd.PersistentFlags().IntVarP(&duration, "duration", "d", pim.DEFAULT_DURATION_MINUTES, "Duration in minutes that the role should be activated for")
	activateCmd.PersistentFlags().StringVar(&reason, "reason", pim.DEFAULT_REASON, "Reason for the activation")
	activateCmd.PersistentFlags().StringVar(&ticketSystem, "ticket-system", "", "Ticket system for the activation")
	activateCmd.PersistentFlags().StringVarP(&ticketNumber, "ticket-number", "T", "", "Ticket number for the activation")
	activateCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Display the resource that would be activated, without requesting the activation")

	activateGroupCmd.PersistentFlags().StringVarP(&pimGovernanceRoleToken, "token", "t", "", "An access token for the PIM 'Entra Roles' and 'Groups' API (required). Consult the README for more information.")
	activateGroupCmd.MarkPersistentFlagRequired("token") //nolint:errcheck

	activateCmd.MarkFlagsOneRequired("name", "prefix")
	activateCmd.MarkFlagsMutuallyExclusive("name", "prefix")
}
