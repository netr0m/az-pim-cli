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

func TestPrintEligibleGovernanceRoles(t *testing.T) {
	PrintEligibleGovernanceRoles(pim.EligibleGovernanceRoleAssignmentsDummyData)
}

func TestGetResourceAssignment(t *testing.T) {
	var sub1role1 *pim.ResourceAssignment = GetResourceAssignment(pim.TEST_DUMMY_SUBSCRIPTION_1_NAME, "", pim.TEST_DUMMY_ROLE_1_NAME, pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, sub1role1, &pim.EligibleResourceAssignmentsDummyData.Value[0], "resulting resource assignment does not match expected assignment")
	var sub1role2 *pim.ResourceAssignment = GetResourceAssignment(pim.TEST_DUMMY_SUBSCRIPTION_1_NAME, "", pim.TEST_DUMMY_ROLE_2_NAME, pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, sub1role2, &pim.EligibleResourceAssignmentsDummyData.Value[1], "resulting resource assignment does not match expected assignment")
	var sub2 *pim.ResourceAssignment = GetResourceAssignment(pim.TEST_DUMMY_SUBSCRIPTION_2_NAME, "", "", pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, sub2, &pim.EligibleResourceAssignmentsDummyData.Value[2], "resulting resource assignment does not match expected assignment")
	assert.Equal(t, sub2.Properties.ExpandedProperties.Scope.DisplayName, pim.TEST_DUMMY_SUBSCRIPTION_2_NAME, "resulting resource assignment scope name does not match expected name")

	var subprefix *pim.ResourceAssignment = GetResourceAssignment("", "azure res", "", pim.EligibleResourceAssignmentsDummyData)
	assert.EqualValues(t, subprefix, &pim.EligibleResourceAssignmentsDummyData.Value[3], "resulting resource assignment does not match expected assignment")
}

func TestGetGovernanceRoleAssignmentAADGroup(t *testing.T) {
	var grp1role1 *pim.GovernanceRoleAssignment = GetGovernanceRoleAssignment(pim.TEST_DUMMY_GROUP_1_NAME, "", pim.TEST_DUMMY_ROLE_1_NAME, pim.EligibleGovernanceRoleAssignmentsDummyData)
	assert.EqualValues(t, grp1role1, &pim.EligibleGovernanceRoleAssignmentsDummyData.Value[0], "resulting governance role assignment does not match expected assignment")
	var grp1role2 *pim.GovernanceRoleAssignment = GetGovernanceRoleAssignment(pim.TEST_DUMMY_GROUP_1_NAME, "", pim.TEST_DUMMY_ROLE_2_NAME, pim.EligibleGovernanceRoleAssignmentsDummyData)
	assert.EqualValues(t, grp1role2, &pim.EligibleGovernanceRoleAssignmentsDummyData.Value[1], "resulting governance role assignment does not match expected assignment")
	assert.Equal(t, grp1role2.RoleDefinition.DisplayName, pim.TEST_DUMMY_ROLE_2_NAME, "resulting governance role assignment role name does not match expected name")
	var grp2role1 *pim.GovernanceRoleAssignment = GetGovernanceRoleAssignment(pim.TEST_DUMMY_GROUP_2_NAME, "", pim.TEST_DUMMY_ROLE_1_NAME, pim.EligibleGovernanceRoleAssignmentsDummyData)
	assert.EqualValues(t, grp2role1, &pim.EligibleGovernanceRoleAssignmentsDummyData.Value[2], "resulting governance role assignment does not match expected assignment")
	assert.Equal(t, grp2role1.RoleDefinition.Resource.DisplayName, pim.TEST_DUMMY_GROUP_2_NAME, "resulting governance role assignment resource name does not match expected name")

	var grpprefix *pim.GovernanceRoleAssignment = GetGovernanceRoleAssignment("", "group", "", pim.EligibleGovernanceRoleAssignmentsDummyData)
	assert.EqualValues(t, grpprefix, &pim.EligibleGovernanceRoleAssignmentsDummyData.Value[0], "resulting governance role assignment does not match expected assignment")
}
