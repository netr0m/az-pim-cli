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

func PrintEligibleResources(resourceAssignments *pim.ResourceAssignmentResponse) {
	var eligibleResources = make(map[string][]string)

	for _, ras := range resourceAssignments.Value {
		resourceName := ras.Properties.ExpandedProperties.Scope.DisplayName
		roleName := ras.Properties.ExpandedProperties.RoleDefinition.DisplayName
		if _, ok := eligibleResources[resourceName]; !ok {
			eligibleResources[resourceName] = []string{}
		}
		eligibleResources[resourceName] = append(eligibleResources[resourceName], roleName)
	}

	for sub, rol := range eligibleResources {
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

func GetResourceAssignment(name string, prefix string, role string, eligibleResourceAssignments *pim.ResourceAssignmentResponse) *pim.ResourceAssignment {
	name = strings.ToLower(name)
	prefix = strings.ToLower(prefix)
	role = strings.ToLower(role)
	for _, eligibleResourceAssignment := range eligibleResourceAssignments.Value {
		var match *pim.ResourceAssignment = nil
		resourceName := strings.ToLower(eligibleResourceAssignment.Properties.ExpandedProperties.Scope.DisplayName)

		if len(prefix) != 0 {
			if strings.HasPrefix(resourceName, prefix) {
				match = &eligibleResourceAssignment
			}
		} else if len(name) != 0 {
			if resourceName == name {
				match = &eligibleResourceAssignment
			}
		}

		if match != nil {
			if role == "" {
				return &eligibleResourceAssignment
			}
			if strings.ToLower(eligibleResourceAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName) == role {
				return &eligibleResourceAssignment
			}
		}
	}

	log.Fatalln("Unable to find a resource assignment matching the parameters.")

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
			if strings.ToLower(eligibleGroupAssignment.RoleDefinition.DisplayName) == role {
				return &eligibleGroupAssignment
			}
		}
	}

	log.Fatalln("Unable to find a group assignment matching the parameters.")

	return nil
}
