/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package pim

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateTime(t *testing.T) {
	now := time.Now().Local()
	currentDate := now.Format("2006-01-02")
	currentTZ := now.Format("-07:00")
	errMsg := "resulting startDateTime does not match expected value"

	dateOnly, _ := parseDateTime("31/12/2024", "")
	timeOnly, _ := parseDateTime("", "13:37")
	dateTime, _ := parseDateTime("31/12/2024", "13:37")

	assert.Equal(t, fmt.Sprintf("2024-12-31T00:00:00%s", currentTZ), dateOnly, errMsg)
	assert.Equal(t, fmt.Sprintf("%sT13:37:00%s", currentDate, currentTZ), timeOnly, errMsg)
	assert.Equal(t, fmt.Sprintf("2024-12-31T13:37:00%s", currentTZ), dateTime, errMsg)
}

func TestGraphRoleEligibilityScheduleInstanceToGovernanceRoleAssignment(t *testing.T) {
	instance := graphRoleEligibilityScheduleInstance{
		Id:               "instance-1",
		PrincipalId:      TEST_DUMMY_PRINCIPAL_ID,
		RoleDefinitionId: "role-def-1",
		DirectoryScopeId: "/",
		RoleDefinition:   &graphRoleDefinition{Id: "role-def-1", DisplayName: "Global Reader"},
		Principal:        &graphPrincipal{Id: TEST_DUMMY_PRINCIPAL_ID, DisplayName: TEST_DUMMY_PRINCIPAL_NAME},
	}

	assignment := instance.toGovernanceRoleAssignment(TEST_DUMMY_PRINCIPAL_ID)

	assert.Equal(t, "instance-1", assignment.Id)
	assert.Equal(t, "/", assignment.ResourceId, "expected ResourceId to hold directoryScopeId for Entra roles")
	assert.Equal(t, "role-def-1", assignment.RoleDefinitionId)
	assert.Equal(t, TEST_DUMMY_PRINCIPAL_ID, assignment.SubjectId)
	assert.Equal(t, "Global Reader", assignment.RoleDefinition.DisplayName)
	assert.Equal(t, "Global Reader", assignment.RoleDefinition.Resource.DisplayName, "expected the Resource grouping to mirror the role name for Entra roles")
	assert.Equal(t, TEST_DUMMY_PRINCIPAL_NAME, assignment.Subject.DisplayName)
}

func TestGraphRoleEligibilityScheduleInstanceToGovernanceRoleAssignmentWithoutExpand(t *testing.T) {
	instance := graphRoleEligibilityScheduleInstance{
		Id:               "instance-1",
		RoleDefinitionId: "role-def-1",
		DirectoryScopeId: "/",
	}

	assignment := instance.toGovernanceRoleAssignment(TEST_DUMMY_PRINCIPAL_ID)

	assert.Equal(t, "role-def-1", assignment.RoleDefinition.DisplayName, "expected a fallback to the raw role definition ID when $expand=roleDefinition is absent")
	assert.Nil(t, assignment.Subject, "expected no Subject when $expand=principal is absent")
}

func TestGraphGroupEligibilityScheduleInstanceToGovernanceRoleAssignment(t *testing.T) {
	instance := graphGroupEligibilityScheduleInstance{
		Id:          "instance-1",
		PrincipalId: TEST_DUMMY_PRINCIPAL_ID,
		GroupId:     TEST_DUMMY_GROUP_1_ID,
		AccessId:    "member",
		Group:       &graphGroup{Id: TEST_DUMMY_GROUP_1_ID, DisplayName: TEST_DUMMY_GROUP_1_NAME},
		Principal:   &graphPrincipal{Id: TEST_DUMMY_PRINCIPAL_ID, DisplayName: TEST_DUMMY_PRINCIPAL_NAME},
	}

	assignment := instance.toGovernanceRoleAssignment(TEST_DUMMY_PRINCIPAL_ID)

	assert.Equal(t, "instance-1", assignment.Id)
	assert.Equal(t, TEST_DUMMY_GROUP_1_ID, assignment.ResourceId)
	assert.Equal(t, "member", assignment.AccessId, "expected AccessId to be kept verbatim for building the activation request")
	assert.Equal(t, "Member", assignment.RoleDefinition.DisplayName, "expected AccessId to be capitalized for display")
	assert.Equal(t, TEST_DUMMY_GROUP_1_NAME, assignment.RoleDefinition.Resource.DisplayName)
	assert.Equal(t, TEST_DUMMY_PRINCIPAL_NAME, assignment.Subject.DisplayName)
}

func TestCreateResourceAssignmentRequest(t *testing.T) {
	resourceAssignment := &EligibleResourceAssignmentsDummyData.Value[0]
	tests := []struct {
		name     string
		scope    string
		expected string
	}{
		{
			name:     "without resource scope",
			scope:    "",
			expected: resourceAssignment.Properties.ExpandedProperties.Scope.Id,
		},
		{
			name:     "with resource scope",
			scope:    fmt.Sprintf("/subscriptions/%s/resourceGroups/rg", TEST_DUMMY_SUBSCRIPTION_1_ID),
			expected: fmt.Sprintf("/subscriptions/%s/resourceGroups/rg", TEST_DUMMY_SUBSCRIPTION_1_ID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope, _ := CreateResourceAssignmentRequestWithScope(TEST_DUMMY_PRINCIPAL_ID, resourceAssignment, tt.scope, 30, "", "", "test", "Test", "1337")

			assert.Equal(t, tt.expected, scope)
		})
	}
}
