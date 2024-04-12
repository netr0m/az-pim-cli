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

func GetRoleAssignment(name interface{}, prefix interface{}, role interface{}, eligibleRoleAssignments *pim.RoleAssignmentResponse) *pim.RoleAssignment {
	if name == nil && prefix == nil {
		log.Fatalf("getSubscriptionId() requires either 'name' or 'prefix' as its input parameter")
	}
	for _, eligibleRoleAssignment := range eligibleRoleAssignments.Value {
		var match *pim.RoleAssignment = nil
		subscriptionName := strings.ToLower(eligibleRoleAssignment.Properties.ExpandedProperties.Scope.DisplayName)

		if prefix, exists := prefix.(string); exists {
			prefix = strings.ToLower(prefix)
			if strings.HasPrefix(subscriptionName, prefix) {
				match = &eligibleRoleAssignment
			}
		} else if name, exists := name.(string); exists {
			name = strings.ToLower(name)
			if subscriptionName == name {
				match = &eligibleRoleAssignment
			}
		}

		if match != nil {
			if role == nil {
				return &eligibleRoleAssignment
			}
			if role, exists := role.(string); exists {
				role = strings.ToLower(role)
				if strings.Contains(strings.ToLower(eligibleRoleAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName), role) {
					return &eligibleRoleAssignment
				}
			}
		}

	}

	log.Fatalln("Unable to find a role assignment matching the parameters.")

	return nil
}
