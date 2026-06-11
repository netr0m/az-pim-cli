/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package utils

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/netr0m/az-pim-cli/pkg/common"
	"github.com/netr0m/az-pim-cli/pkg/pim"
)

func PrintEligibleResources(resourceAssignments *pim.ResourceAssignmentResponse) {
	var eligibleResources = make(map[string][]string)

	for _, ras := range resourceAssignments.Value {
		slog.Debug(ras.Debug())
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

func PrintEligibleGroups(groupAssignments *pim.GraphGroupEligibilityResponse) {
	var eligibleGroups = make(map[string][]string)

	for i := range groupAssignments.Value {
		instance := &groupAssignments.Value[i]
		slog.Debug(instance.Debug())
		groupName := "<unknown>"
		if instance.Group != nil {
			groupName = instance.Group.DisplayName
		}
		eligibleGroups[groupName] = append(eligibleGroups[groupName], instance.AccessId)
	}

	for groupName, accessIds := range eligibleGroups {
		fmt.Printf("== %s ==\n", groupName)
		for _, accessId := range accessIds {
			fmt.Printf("\t - %s\n", accessId)
		}
	}
}

func PrintEligibleRoles(roleAssignments *pim.GraphRoleEligibilityResponse) {
	fmt.Println("== Entra roles ==")
	for i := range roleAssignments.Value {
		instance := &roleAssignments.Value[i]
		slog.Debug(instance.Debug())
		roleName := "<unknown>"
		if instance.RoleDefinition != nil {
			roleName = instance.RoleDefinition.DisplayName
		}
		fmt.Printf("\t - %s\n", roleName)
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

	var _error = common.Error{
		Operation: "GetResourceAssignment",
		Message:   "Unable to find a resource assignment matching the parameters",
		Status:    "404",
	}
	slog.Error(_error.Error())
	os.Exit(1)

	return nil
}

func GetEligibleGroupAssignment(name string, prefix string, role string, eligibleGroupAssignments *pim.GraphGroupEligibilityResponse) *pim.GraphGroupEligibilityInstance {
	name = strings.ToLower(name)
	prefix = strings.ToLower(prefix)
	role = strings.ToLower(role)
	for i := range eligibleGroupAssignments.Value {
		instance := &eligibleGroupAssignments.Value[i]
		groupName := ""
		if instance.Group != nil {
			groupName = strings.ToLower(instance.Group.DisplayName)
		}

		var matched bool
		if len(prefix) != 0 {
			matched = strings.HasPrefix(groupName, prefix)
		} else if len(name) != 0 {
			matched = groupName == name
		}

		if matched {
			// For groups, the "role" is the access type (member/owner)
			if role == "" || strings.ToLower(instance.AccessId) == role {
				return instance
			}
		}
	}

	var _error = common.Error{
		Operation: "GetEligibleGroupAssignment",
		Message:   "Unable to find a group assignment matching the parameters",
		Status:    "404",
	}
	slog.Error(_error.Error())
	os.Exit(1)

	return nil
}

func GetEligibleRoleAssignment(name string, prefix string, role string, eligibleRoleAssignments *pim.GraphRoleEligibilityResponse) *pim.GraphRoleEligibilityInstance {
	name = strings.ToLower(name)
	prefix = strings.ToLower(prefix)
	role = strings.ToLower(role)
	for i := range eligibleRoleAssignments.Value {
		instance := &eligibleRoleAssignments.Value[i]
		roleName := ""
		if instance.RoleDefinition != nil {
			roleName = strings.ToLower(instance.RoleDefinition.DisplayName)
		}

		var matched bool
		if len(prefix) != 0 {
			matched = strings.HasPrefix(roleName, prefix)
		} else if len(name) != 0 {
			matched = roleName == name
		}

		if matched {
			// For Entra roles the container and the role are the same; '--role'
			// (if provided) is matched against the same role display name.
			if role == "" || roleName == role {
				return instance
			}
		}
	}

	var _error = common.Error{
		Operation: "GetEligibleRoleAssignment",
		Message:   "Unable to find a role assignment matching the parameters",
		Status:    "404",
	}
	slog.Error(_error.Error())
	os.Exit(1)

	return nil
}
