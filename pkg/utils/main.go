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
	name = strings.ToLower(name)
	prefix = strings.ToLower(prefix)
	role = strings.ToLower(role)
	for _, eligibleRoleAssignment := range eligibleRoleAssignments.Value {
		var match *pim.RoleAssignment = nil
		subscriptionName := strings.ToLower(eligibleRoleAssignment.Properties.ExpandedProperties.Scope.DisplayName)

		if len(prefix) != 0 {
			if strings.HasPrefix(subscriptionName, prefix) {
				match = &eligibleRoleAssignment
			}
		} else if len(name) != 0 {
			if subscriptionName == name {
				match = &eligibleRoleAssignment
			}
		}

		if match != nil {
			if role == "" {
				return &eligibleRoleAssignment
			}
			if strings.Contains(strings.ToLower(eligibleRoleAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName), role) {
				return &eligibleRoleAssignment
			}
		}
	}

	log.Fatalln("Unable to find a role assignment matching the parameters.")

	return nil
}

func GetGroupAssignment(name string, prefix string, role string, eligibleGroupAssignments *pim.GroupAssignmentResponse) *pim.GroupAssignment {
	name = strings.ToLower(name)
	prefix = strings.ToLower(prefix)
	role = strings.ToLower(role)
	for _, eligibleGroupAssignment := range eligibleGroupAssignments.Value {
		var match *pim.GroupAssignment = nil
		currentGroupName := strings.ToLower(eligibleGroupAssignment.RoleDefinition.Resource.DisplayName)

		if len(prefix) != 0 {
			if strings.HasPrefix(currentGroupName, prefix) {
				match = &eligibleGroupAssignment // #nosec G601 false positive with go >= v1.22
			}
		} else if len(name) != 0 {
			if currentGroupName == name {
				match = &eligibleGroupAssignment // #nosec G601 false positive with go >= v1.22
			}
		}

		if match != nil {
			if role == "" {
				return &eligibleGroupAssignment
			}
			if strings.Contains(strings.ToLower(eligibleGroupAssignment.RoleDefinition.DisplayName), role) {
				return &eligibleGroupAssignment
			}
		}
	}

	log.Fatalln("Unable to find a group assignment matching the parameters.")

	return nil
}
