/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package utils

import (
	"testing"

	"github.com/netr0m/az-pim-cli/pkg/pim"
	"github.com/stretchr/testify/assert"
)

func TestPrintEligibleResources(t *testing.T) {
	PrintEligibleResources(pim.EligibleResourceAssignmentsDummyData)
}

func TestPrintEligibleGroups(t *testing.T) {
	PrintEligibleGroups(pim.EligibleGroupAssignmentsDummyData)
}

func TestPrintEligibleRoles(t *testing.T) {
	PrintEligibleRoles(pim.EligibleRoleAssignmentsDummyData)
}

func TestGetResourceAssignment(t *testing.T) {
	var sub1role1 = GetResourceAssignment(pim.TEST_DUMMY_SUBSCRIPTION_1_NAME, "", pim.TEST_DUMMY_ROLE_1_NAME, pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, sub1role1, &pim.EligibleResourceAssignmentsDummyData.Value[0], "resulting resource assignment does not match expected assignment")
	var sub1role2 = GetResourceAssignment(pim.TEST_DUMMY_SUBSCRIPTION_1_NAME, "", pim.TEST_DUMMY_ROLE_2_NAME, pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, sub1role2, &pim.EligibleResourceAssignmentsDummyData.Value[1], "resulting resource assignment does not match expected assignment")
	var sub2 = GetResourceAssignment(pim.TEST_DUMMY_SUBSCRIPTION_2_NAME, "", "", pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, sub2, &pim.EligibleResourceAssignmentsDummyData.Value[2], "resulting resource assignment does not match expected assignment")
	assert.Equal(t, sub2.Properties.ExpandedProperties.Scope.DisplayName, pim.TEST_DUMMY_SUBSCRIPTION_2_NAME, "resulting resource assignment scope name does not match expected name")

	var subprefix = GetResourceAssignment("", "azure res", "", pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, subprefix, &pim.EligibleResourceAssignmentsDummyData.Value[3], "resulting resource assignment does not match expected assignment")
}

func TestGetEligibleGroupAssignment(t *testing.T) {
	// Group 1 has both a 'member' (Value[0]) and an 'owner' (Value[1]) eligibility
	var grp1member = GetEligibleGroupAssignment(pim.TEST_DUMMY_GROUP_1_NAME, "", "member", pim.EligibleGroupAssignmentsDummyData)
	assert.EqualValues(t, grp1member, &pim.EligibleGroupAssignmentsDummyData.Value[0], "resulting group assignment does not match expected assignment")
	var grp1owner = GetEligibleGroupAssignment(pim.TEST_DUMMY_GROUP_1_NAME, "", "owner", pim.EligibleGroupAssignmentsDummyData)
	assert.EqualValues(t, grp1owner, &pim.EligibleGroupAssignmentsDummyData.Value[1], "resulting group assignment does not match expected assignment")
	assert.Equal(t, grp1owner.AccessId, "owner", "resulting group assignment accessId does not match expected value")
	// Group 2 has a single 'member' eligibility; no role filter required
	var grp2 = GetEligibleGroupAssignment(pim.TEST_DUMMY_GROUP_2_NAME, "", "", pim.EligibleGroupAssignmentsDummyData)
	assert.EqualValues(t, grp2, &pim.EligibleGroupAssignmentsDummyData.Value[2], "resulting group assignment does not match expected assignment")
	assert.Equal(t, grp2.Group.DisplayName, pim.TEST_DUMMY_GROUP_2_NAME, "resulting group assignment group name does not match expected name")

	var grpprefix = GetEligibleGroupAssignment("", "group", "", pim.EligibleGroupAssignmentsDummyData)
	assert.EqualValues(t, grpprefix, &pim.EligibleGroupAssignmentsDummyData.Value[0], "resulting group assignment does not match expected assignment")
}

func TestGetEligibleRoleAssignment(t *testing.T) {
	var role1 = GetEligibleRoleAssignment(pim.TEST_DUMMY_ROLE_1_NAME, "", "", pim.EligibleRoleAssignmentsDummyData)
	assert.EqualValues(t, role1, &pim.EligibleRoleAssignmentsDummyData.Value[0], "resulting role assignment does not match expected assignment")
	var role2 = GetEligibleRoleAssignment(pim.TEST_DUMMY_ROLE_2_NAME, "", pim.TEST_DUMMY_ROLE_2_NAME, pim.EligibleRoleAssignmentsDummyData)
	assert.EqualValues(t, role2, &pim.EligibleRoleAssignmentsDummyData.Value[1], "resulting role assignment does not match expected assignment")
	assert.Equal(t, role2.RoleDefinition.DisplayName, pim.TEST_DUMMY_ROLE_2_NAME, "resulting role assignment role name does not match expected name")

	var roleprefix = GetEligibleRoleAssignment("", "role 1", "", pim.EligibleRoleAssignmentsDummyData)
	assert.EqualValues(t, roleprefix, &pim.EligibleRoleAssignmentsDummyData.Value[0], "resulting role assignment does not match expected assignment")
}
