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

func PrintEligibleGovernanceRoles(governanceRoleAssignments *pim.GovernanceRoleAssignmentResponse) {
	var eligibleGovernanceRoles = make(map[string][]string)

	for _, governanceRoleAssignment := range governanceRoleAssignments.Value {
		governanceRoleName := governanceRoleAssignment.RoleDefinition.Resource.DisplayName
		roleName := governanceRoleAssignment.RoleDefinition.DisplayName
		if _, ok := eligibleGovernanceRoles[governanceRoleName]; !ok {
			eligibleGovernanceRoles[governanceRoleName] = []string{}
		}
		eligibleGovernanceRoles[governanceRoleName] = append(eligibleGovernanceRoles[governanceRoleName], roleName)
	}

	for govRole, rol := range eligibleGovernanceRoles {
		fmt.Printf("== %s ==\n", govRole)
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

func GetGovernanceRoleAssignment(name string, prefix string, role string, eligibleGovernanceRoleAssignments *pim.GovernanceRoleAssignmentResponse) *pim.GovernanceRoleAssignment {
	name = strings.ToLower(name)
	prefix = strings.ToLower(prefix)
	role = strings.ToLower(role)
	for _, eligibleGovernanceRoleAssignment := range eligibleGovernanceRoleAssignments.Value {
		var match *pim.GovernanceRoleAssignment = nil
		currentGovernanceRoleName := strings.ToLower(eligibleGovernanceRoleAssignment.RoleDefinition.Resource.DisplayName)

		if len(prefix) != 0 {
			if strings.HasPrefix(currentGovernanceRoleName, prefix) {
				match = &eligibleGovernanceRoleAssignment // #nosec G601 false positive with go >= v1.22
			}
		} else if len(name) != 0 {
			if currentGovernanceRoleName == name {
				match = &eligibleGovernanceRoleAssignment // #nosec G601 false positive with go >= v1.22
			}
		}

		if match != nil {
			if role == "" {
				return &eligibleGovernanceRoleAssignment
			}
			if strings.ToLower(eligibleGovernanceRoleAssignment.RoleDefinition.DisplayName) == role {
				return &eligibleGovernanceRoleAssignment
			}
		}
	}

	log.Fatalln("Unable to find a role assignment matching the parameters.")

	return nil
}
