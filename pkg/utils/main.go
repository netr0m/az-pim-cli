/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/netr0m/az-pim-cli/pkg/pim"
)

func PrintEligibleRoles(roleEligibilityScheduleInstances *pim.RoleAssignmentResponse) {
	var eligibleRoles = make(map[string][]string)

	for _, ras := range roleEligibilityScheduleInstances.Value {
		subscriptionName := ras.Properties.ExpandedProperties.Scope.DisplayName
		roleName := ras.Properties.ExpandedProperties.RoleDefinition.DisplayName
		if _, ok := eligibleRoles[subscriptionName]; !ok {
			eligibleRoles[subscriptionName] = []string{}
		}
		eligibleRoles[subscriptionName] = append(eligibleRoles[subscriptionName], roleName)
	}

	for sub, rol := range eligibleRoles {
		fmt.Printf("== %s ==\n", sub)
		for role := range rol {
			fmt.Printf("\t - %s\n", rol[role])
		}
	}
}

func PrintEligibleGroups(groupAssignments *pim.GroupAssignmentResponse) {
	var eligibleGroups = make(map[string][]string)

	for _, groupAssignment := range groupAssignments.Value {
		groupName := groupAssignment.RoleDefinition.Resource.DisplayName
		roleName := groupAssignment.RoleDefinition.DisplayName
		if _, ok := eligibleGroups[groupName]; !ok {
			eligibleGroups[groupName] = []string{}
		}
		eligibleGroups[groupName] = append(eligibleGroups[groupName], roleName)
	}

	for grp, rol := range eligibleGroups {
		fmt.Printf("== %s ==\n", grp)
		for role := range rol {
			fmt.Printf("\t - %s\n", rol[role])
		}
	}
}

func GetRoleAssignment(name string, prefix string, role string, eligibleRoleAssignments *pim.RoleAssignmentResponse) *pim.RoleAssignment {
	for _, eligibleRoleAssignment := range eligibleRoleAssignments.Value {
		var match *pim.RoleAssignment = nil
		subscriptionName := strings.ToLower(eligibleRoleAssignment.Properties.ExpandedProperties.Scope.DisplayName)

		if len(prefix) != 0 {
			prefix = strings.ToLower(prefix)
			if strings.HasPrefix(subscriptionName, prefix) {
				match = &eligibleRoleAssignment
			}
		} else if len(name) != 0 {
			name = strings.ToLower(name)
			if subscriptionName == name {
				match = &eligibleRoleAssignment
			}
		}

		if match != nil {
			if role == "" {
				return &eligibleRoleAssignment
			}
			role = strings.ToLower(role)
			if strings.Contains(strings.ToLower(eligibleRoleAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName), role) {
				return &eligibleRoleAssignment
			}
		}
	}

	log.Fatalln("Unable to find a role assignment matching the parameters.")

	return nil
}
