/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"os"

	"log/slog"

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
		token := pim.GetAccessToken(pim.AZ_PIM_SCOPE, pim.AzureClient{})
		subjectId := pim.GetUserInfo(token).ObjectId

		eligibleResourceAssignments := pim.GetEligibleResourceAssignments(token, pim.AzureClient{})
		resourceAssignment := utils.GetResourceAssignment(name, prefix, roleName, eligibleResourceAssignments)

		slog.Info(
			"Requesting activation",
			"role", resourceAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName,
			"scope", resourceAssignment.Properties.ExpandedProperties.Scope.DisplayName,
			"reason", reason,
			"ticketNumber", ticketNumber,
			"ticketSystem", ticketSystem,
		)

		if dryRun {
			slog.Warn("Skipping activation due to '--dry-run'")
			os.Exit(0)
		}
		requestResponse := pim.RequestResourceAssignment(subjectId, resourceAssignment, duration, reason, ticketSystem, ticketNumber, token)
		slog.Info(
			"Request completed",
			"role", resourceAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName,
			"scope", resourceAssignment.Properties.ExpandedProperties.Scope.DisplayName,
			"status", requestResponse.Properties.Status,
		)
	},
}

var activateGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"g", "grp", "groups"},
	Short:   "Sends a request to Azure PIM to activate the given group",
	Run: func(cmd *cobra.Command, args []string) {
		subjectId := pim.GetUserInfo(pimGovernanceRoleToken).ObjectId

		eligibleGroupAssignments := pim.GetEligibleGovernanceRoleAssignments(pim.ROLE_TYPE_AAD_GROUPS, subjectId, pimGovernanceRoleToken, pim.AzureClient{})
		groupAssignment := utils.GetGovernanceRoleAssignment(name, prefix, roleName, eligibleGroupAssignments)

		slog.Info(
			"Requesting activation",
			"role", groupAssignment.RoleDefinition.DisplayName,
			"scope", groupAssignment.RoleDefinition.Resource.DisplayName,
			"reason", reason,
			"ticketNumber", ticketNumber,
			"ticketSystem", ticketSystem,
		)

		if dryRun {
			slog.Warn("Skipping activation due to '--dry-run'")
			os.Exit(0)
		}
		requestResponse := pim.RequestGovernanceRoleAssignment(subjectId, pim.ROLE_TYPE_AAD_GROUPS, groupAssignment, duration, reason, ticketSystem, ticketNumber, pimGovernanceRoleToken)
		slog.Info(
			"Request completed",
			"role", groupAssignment.RoleDefinition.DisplayName,
			"scope", groupAssignment.RoleDefinition.Resource.DisplayName,
			"status", requestResponse.AssignmentState,
		)
	},
}

var activateEntraRoleCmd = &cobra.Command{
	Use:     "role",
	Aliases: []string{"rl", "role", "roles"},
	Short:   "Sends a request to Azure PIM to activate the given Entra role",
	Run: func(cmd *cobra.Command, args []string) {
		subjectId := pim.GetUserInfo(pimGovernanceRoleToken).ObjectId

		eligibleEntraRoleAssignments := pim.GetEligibleGovernanceRoleAssignments(pim.ROLE_TYPE_ENTRA_ROLES, subjectId, pimGovernanceRoleToken, pim.AzureClient{})
		entraRoleAssignment := utils.GetGovernanceRoleAssignment(name, prefix, roleName, eligibleEntraRoleAssignments)

		slog.Info(
			"Requesting activation",
			"role", entraRoleAssignment.RoleDefinition.DisplayName,
			"scope", entraRoleAssignment.RoleDefinition.Resource.DisplayName,
			"reason", reason,
			"ticketNumber", ticketNumber,
			"ticketSystem", ticketSystem,
		)

		if dryRun {
			slog.Warn("Skipping activation due to '--dry-run'")
			os.Exit(0)
		}
		requestResponse := pim.RequestGovernanceRoleAssignment(subjectId, pim.ROLE_TYPE_ENTRA_ROLES, entraRoleAssignment, duration, reason, ticketSystem, ticketNumber, pimGovernanceRoleToken)
		slog.Info(
			"Request completed",
			"role", entraRoleAssignment.RoleDefinition.DisplayName,
			"scope", entraRoleAssignment.RoleDefinition.Resource.DisplayName,
			"status", requestResponse.AssignmentState,
		)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.AddCommand(activateResourceCmd)
	activateCmd.AddCommand(activateGroupCmd)
	activateCmd.AddCommand(activateEntraRoleCmd)

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

	activateEntraRoleCmd.PersistentFlags().StringVarP(&pimGovernanceRoleToken, "token", "t", "", "An access token for the PIM 'Entra Roles' and 'Groups' API (required). Consult the README for more information.")
	activateEntraRoleCmd.MarkPersistentFlagRequired("token") //nolint:errcheck

	activateCmd.MarkFlagsOneRequired("name", "prefix")
	activateCmd.MarkFlagsMutuallyExclusive("name", "prefix")
}
