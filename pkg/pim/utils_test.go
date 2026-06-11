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
	// The timezone offset must be derived from the parsed date itself, not from
	// "now", since daylight saving time means the offset can differ between the
	// current date and 31/12/2024.
	dec2024TZ := time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local).Format("-07:00")
	errMsg := "resulting startDateTime does not match expected value"

	dateOnly, _ := parseDateTime("31/12/2024", "")
	timeOnly, _ := parseDateTime("", "13:37")
	dateTime, _ := parseDateTime("31/12/2024", "13:37")

	assert.Equal(t, fmt.Sprintf("2024-12-31T00:00:00%s", dec2024TZ), dateOnly, errMsg)
	assert.Equal(t, fmt.Sprintf("%sT13:37:00%s", currentDate, currentTZ), timeOnly, errMsg)
	assert.Equal(t, fmt.Sprintf("2024-12-31T13:37:00%s", dec2024TZ), dateTime, errMsg)
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

func TestCreateGraphScheduleInfo(t *testing.T) {
	tests := []struct {
		name          string
		startDate     string
		startTime     string
		expectStart   bool
		expectedStart string
	}{
		{name: "immediate (no start)", startDate: "", startTime: "", expectStart: false},
		{name: "explicit start date", startDate: "31/12/2024", startTime: "13:37", expectStart: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduleInfo := CreateGraphScheduleInfo(30, tt.startDate, tt.startTime)
			assert.Equal(t, "PT30M", scheduleInfo.Expiration.Duration)
			assert.Equal(t, GRAPH_EXPIRATION_AFTER_DURATION, scheduleInfo.Expiration.Type)
			if tt.expectStart {
				assert.NotNil(t, scheduleInfo.StartDateTime)
			} else {
				assert.Nil(t, scheduleInfo.StartDateTime)
			}
		})
	}
}

func TestCreateGraphGroupAssignmentRequest(t *testing.T) {
	groupAssignment := &EligibleGroupAssignmentsDummyData.Value[1] // Group 1 / owner
	request := CreateGraphGroupAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, groupAssignment, DEFAULT_DURATION_MINUTES, "", "", DEFAULT_REASON, "Test", "1337")

	assert.Equal(t, GRAPH_ACTION_SELF_ACTIVATE, request.Action)
	assert.Equal(t, "owner", request.AccessId)
	assert.Equal(t, TEST_DUMMY_GROUP_1_ID, request.GroupId)
	assert.Equal(t, TEST_DUMMY_PRINCIPAL_ID, request.PrincipalId)
	assert.Equal(t, DEFAULT_REASON, request.Justification)
	assert.Equal(t, fmt.Sprintf("PT%dM", DEFAULT_DURATION_MINUTES), request.ScheduleInfo.Expiration.Duration)
}

func TestCreateGraphRoleAssignmentRequest(t *testing.T) {
	roleAssignment := &EligibleRoleAssignmentsDummyData.Value[0]
	request := CreateGraphRoleAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, roleAssignment, DEFAULT_DURATION_MINUTES, "", "", DEFAULT_REASON, "Test", "1337")

	assert.Equal(t, GRAPH_ACTION_SELF_ACTIVATE, request.Action)
	assert.Equal(t, TEST_DUMMY_ROLE_1_DEFINITION_ID, request.RoleDefinitionId)
	assert.Equal(t, GRAPH_DEFAULT_DIRECTORY_SCOPE, request.DirectoryScopeId)
	assert.Equal(t, TEST_DUMMY_PRINCIPAL_ID, request.PrincipalId)
	assert.Equal(t, fmt.Sprintf("PT%dM", DEFAULT_DURATION_MINUTES), request.ScheduleInfo.Expiration.Duration)
}
