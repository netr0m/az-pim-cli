/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package pim

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct{ mock.Mock }

func newMockClient() *mockClient { return &mockClient{} }

func (m *mockClient) GetAccessToken(scope string) string {
	args := m.Called(scope)
	return args.String(0)
}

func TestGetAccessToken(t *testing.T) {
	m := newMockClient()

	m.On("GetAccessToken", ARM_GLOBAL_BASE_URL).Return(TEST_DUMMY_JWT)

	token := GetAccessToken(ARM_GLOBAL_BASE_URL, m)

	if !strings.HasPrefix(token, "ey") {
		t.Errorf("expected token to start with 'ey', got %s", token)
	}

	m.AssertCalled(t, "GetAccessToken", ARM_GLOBAL_BASE_URL)
}

func TestGetUserInfo(t *testing.T) {
	m := newMockClient()

	m.On("GetAccessToken", ARM_GLOBAL_BASE_URL).Return(TEST_DUMMY_JWT)

	token := GetAccessToken(ARM_GLOBAL_BASE_URL, m)
	userInfo := GetUserInfo(token)

	if userInfo.Email != TEST_DUMMY_PRINCIPAL_EMAIL {
		t.Errorf("unexpected value for userInfo.Email, got %s", userInfo.Email)
	}
}

func (m *mockClient) GetEligibleResourceAssignments(token string) *ResourceAssignmentResponse {
	args := m.Called(token)
	return args.Get(0).(*ResourceAssignmentResponse)
}

func TestGetEligibleResourceAssignments(t *testing.T) {
	m := newMockClient()

	m.On("GetEligibleResourceAssignments", TEST_DUMMY_JWT).Return(EligibleResourceAssignmentsDummyData)

	eligibleResourceAssignments := GetEligibleResourceAssignments(TEST_DUMMY_JWT, m)

	if len(eligibleResourceAssignments.Value) != 4 {
		t.Errorf("expected 4 eligible resource assignments, got %v", len(eligibleResourceAssignments.Value))
	}
	for _, governanceRole := range eligibleResourceAssignments.Value {
		_principalId := governanceRole.Properties.ExpandedProperties.Principal.Id
		if _principalId != TEST_DUMMY_PRINCIPAL_ID {
			t.Errorf("expected resource Properties.ExpandedProperties.Principal.Id to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, _principalId)
		}
	}
	// Check resource name
	_resourceName := eligibleResourceAssignments.Value[1].Properties.ExpandedProperties.Scope.DisplayName
	if _resourceName != TEST_DUMMY_SUBSCRIPTION_1_NAME {
		t.Errorf("expected resource Properties.ExpandedProperties.Scope.DisplayName to be %s, got %s", TEST_DUMMY_SUBSCRIPTION_1_NAME, _resourceName)
	}
	// Check role name
	_roleName := eligibleResourceAssignments.Value[2].Properties.ExpandedProperties.RoleDefinition.DisplayName
	if _roleName != TEST_DUMMY_ROLE_1_NAME {
		t.Errorf("expected resource Properties.ExpandedProperties.RoleDefinition.DisplayName to be %s, got %s", TEST_DUMMY_ROLE_1_NAME, _roleName)
	}
}

func (m *mockClient) GetEligibleGroupAssignments(principalId string, token string) *GraphGroupEligibilityResponse {
	args := m.Called(principalId, token)
	return args.Get(0).(*GraphGroupEligibilityResponse)
}

func TestGetEligibleGroupAssignments(t *testing.T) {
	m := newMockClient()

	m.On("GetEligibleGroupAssignments", TEST_DUMMY_PRINCIPAL_ID, TEST_DUMMY_JWT).Return(EligibleGroupAssignmentsDummyData)

	eligibleGroupAssignments := GetEligibleGroupAssignments(TEST_DUMMY_PRINCIPAL_ID, TEST_DUMMY_JWT, m)

	if len(eligibleGroupAssignments.Value) != 3 {
		t.Errorf("expected 3 eligible group assignments, got %v", len(eligibleGroupAssignments.Value))
	}
	for _, groupAssignment := range eligibleGroupAssignments.Value {
		if groupAssignment.PrincipalId != TEST_DUMMY_PRINCIPAL_ID {
			t.Errorf("expected group assignment PrincipalId to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, groupAssignment.PrincipalId)
		}
	}
	// Check group name
	_groupName := eligibleGroupAssignments.Value[1].Group.DisplayName
	if _groupName != TEST_DUMMY_GROUP_1_NAME {
		t.Errorf("expected group assignment Group.DisplayName to be %s, got %s", TEST_DUMMY_GROUP_1_NAME, _groupName)
	}
	// Check access id
	_accessId := eligibleGroupAssignments.Value[1].AccessId
	if _accessId != "owner" {
		t.Errorf("expected group assignment AccessId to be %s, got %s", "owner", _accessId)
	}
}

func (m *mockClient) GetEligibleRoleAssignments(principalId string, token string) *GraphRoleEligibilityResponse {
	args := m.Called(principalId, token)
	return args.Get(0).(*GraphRoleEligibilityResponse)
}

func TestGetEligibleRoleAssignments(t *testing.T) {
	m := newMockClient()

	m.On("GetEligibleRoleAssignments", TEST_DUMMY_PRINCIPAL_ID, TEST_DUMMY_JWT).Return(EligibleRoleAssignmentsDummyData)

	eligibleRoleAssignments := GetEligibleRoleAssignments(TEST_DUMMY_PRINCIPAL_ID, TEST_DUMMY_JWT, m)

	if len(eligibleRoleAssignments.Value) != 2 {
		t.Errorf("expected 2 eligible role assignments, got %v", len(eligibleRoleAssignments.Value))
	}
	for _, roleAssignment := range eligibleRoleAssignments.Value {
		if roleAssignment.PrincipalId != TEST_DUMMY_PRINCIPAL_ID {
			t.Errorf("expected role assignment PrincipalId to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, roleAssignment.PrincipalId)
		}
	}
	// Check role name
	_roleName := eligibleRoleAssignments.Value[0].RoleDefinition.DisplayName
	if _roleName != TEST_DUMMY_ROLE_1_NAME {
		t.Errorf("expected role assignment RoleDefinition.DisplayName to be %s, got %s", TEST_DUMMY_ROLE_1_NAME, _roleName)
	}
}

func (m *mockClient) ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) bool {
	args := m.Called(scope, resourceAssignmentRequest, token)
	return args.Bool(0)
}

func TestValidateResourceAssignmentRequest(t *testing.T) {
	m := newMockClient()

	resourceAssignment := &EligibleResourceAssignmentsDummyData.Value[0]
	scope, resourceAssignmentRequest := CreateResourceAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, resourceAssignment, 30, "", "", "test", "Test", "1337")

	m.On("ValidateResourceAssignmentRequest", scope, resourceAssignmentRequest, TEST_DUMMY_JWT).Return(true)

	isValid := ValidateResourceAssignmentRequest(scope, resourceAssignmentRequest, TEST_DUMMY_JWT, m)

	if !isValid {
		t.Errorf("expected resource assignment request validation to be successful, got %v", isValid)
	}
}

func (m *mockClient) RequestResourceAssignment(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) *ResourceAssignmentRequestResponse {
	args := m.Called(scope, resourceAssignmentRequest, token)
	return args.Get(0).(*ResourceAssignmentRequestResponse)
}

func TestRequestResourceAssignment(t *testing.T) {
	m := newMockClient()

	resourceAssignment := &EligibleResourceAssignmentsDummyData.Value[0]
	scope, resourceAssignmentRequest := CreateResourceAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, resourceAssignment, DEFAULT_DURATION_MINUTES, "", "", DEFAULT_REASON, "Test", "1337")
	resourceAssignmentRequestResponse := &ResourceAssignmentRequestResponse{
		Id:   resourceAssignment.Id,
		Name: resourceAssignment.Name,
		Type: resourceAssignment.Type,
		Properties: &ResourceAssignmentValidationProperties{
			Scope:              scope,
			PrincipalId:        resourceAssignmentRequest.Properties.PrincipalId,
			Status:             "Active",
			ScheduleInfo:       resourceAssignmentRequest.Properties.ScheduleInfo,
			Justification:      DEFAULT_REASON,
			TicketInfo:         resourceAssignmentRequest.Properties.TicketInfo,
			RoleDefinitionId:   resourceAssignmentRequest.Properties.RoleDefinitionId,
			ExpandedProperties: resourceAssignment.Properties.ExpandedProperties,
		},
	}

	m.On("RequestResourceAssignment", scope, resourceAssignmentRequest, TEST_DUMMY_JWT).Return(resourceAssignmentRequestResponse)

	requestResponse := RequestResourceAssignment(scope, resourceAssignmentRequest, TEST_DUMMY_JWT, m)
	expectedDuration := fmt.Sprintf("PT%dM", DEFAULT_DURATION_MINUTES)

	assert.Equal(t, requestResponse.Properties.Justification, DEFAULT_REASON, "expected resource assignment request justification to be %s, got %s", DEFAULT_REASON, requestResponse.Properties.Justification)
	assert.Equal(t, requestResponse.Properties.PrincipalId, TEST_DUMMY_PRINCIPAL_ID, "expected resource assignment request principal ID to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, requestResponse.Properties.PrincipalId)
	assert.Equal(t, requestResponse.Properties.Status, "Active", "expected resource assignment request status to be %s, got %s", "Active", requestResponse.Properties.Status)
	assert.Equal(t, requestResponse.Properties.ScheduleInfo.Expiration.Duration, expectedDuration, "expected resource assignment request expiration duration to be %s, got %s", expectedDuration, requestResponse.Properties.Status)
}

func (m *mockClient) RequestGroupAssignment(groupAssignmentRequest *GraphGroupAssignmentRequest, token string) *GraphAssignmentScheduleRequest {
	args := m.Called(groupAssignmentRequest, token)
	return args.Get(0).(*GraphAssignmentScheduleRequest)
}

func TestRequestGroupAssignment(t *testing.T) {
	m := newMockClient()

	groupAssignment := &EligibleGroupAssignmentsDummyData.Value[0]
	groupAssignmentRequest := CreateGraphGroupAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, groupAssignment, DEFAULT_DURATION_MINUTES, "", "", DEFAULT_REASON, "Test", "1337")
	groupAssignmentResponse := &GraphAssignmentScheduleRequest{
		Id:     groupAssignment.Id,
		Status: StatusProvisioned,
		Action: GRAPH_ACTION_SELF_ACTIVATE,
	}

	m.On("RequestGroupAssignment", groupAssignmentRequest, TEST_DUMMY_JWT).Return(groupAssignmentResponse)

	requestResponse := RequestGroupAssignment(groupAssignmentRequest, TEST_DUMMY_JWT, m)
	expectedDuration := fmt.Sprintf("PT%dM", DEFAULT_DURATION_MINUTES)

	assert.Equal(t, GRAPH_ACTION_SELF_ACTIVATE, groupAssignmentRequest.Action, "expected group assignment request action to be %s, got %s", GRAPH_ACTION_SELF_ACTIVATE, groupAssignmentRequest.Action)
	assert.Equal(t, groupAssignment.AccessId, groupAssignmentRequest.AccessId, "expected group assignment request accessId to be %s, got %s", groupAssignment.AccessId, groupAssignmentRequest.AccessId)
	assert.Equal(t, TEST_DUMMY_PRINCIPAL_ID, groupAssignmentRequest.PrincipalId, "expected group assignment request principalId to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, groupAssignmentRequest.PrincipalId)
	assert.Equal(t, expectedDuration, groupAssignmentRequest.ScheduleInfo.Expiration.Duration, "expected group assignment request duration to be %s, got %s", expectedDuration, groupAssignmentRequest.ScheduleInfo.Expiration.Duration)
	assert.Equal(t, StatusProvisioned, requestResponse.Status, "expected group assignment response status to be %s, got %s", StatusProvisioned, requestResponse.Status)
}

func (m *mockClient) RequestRoleAssignment(roleAssignmentRequest *GraphRoleAssignmentRequest, token string) *GraphAssignmentScheduleRequest {
	args := m.Called(roleAssignmentRequest, token)
	return args.Get(0).(*GraphAssignmentScheduleRequest)
}

func TestRequestRoleAssignment(t *testing.T) {
	m := newMockClient()

	roleAssignment := &EligibleRoleAssignmentsDummyData.Value[0]
	roleAssignmentRequest := CreateGraphRoleAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, roleAssignment, DEFAULT_DURATION_MINUTES, "", "", DEFAULT_REASON, "Test", "1337")
	roleAssignmentResponse := &GraphAssignmentScheduleRequest{
		Id:     roleAssignment.Id,
		Status: StatusProvisioned,
		Action: GRAPH_ACTION_SELF_ACTIVATE,
	}

	m.On("RequestRoleAssignment", roleAssignmentRequest, TEST_DUMMY_JWT).Return(roleAssignmentResponse)

	requestResponse := RequestRoleAssignment(roleAssignmentRequest, TEST_DUMMY_JWT, m)
	expectedDuration := fmt.Sprintf("PT%dM", DEFAULT_DURATION_MINUTES)

	assert.Equal(t, GRAPH_ACTION_SELF_ACTIVATE, roleAssignmentRequest.Action, "expected role assignment request action to be %s, got %s", GRAPH_ACTION_SELF_ACTIVATE, roleAssignmentRequest.Action)
	assert.Equal(t, GRAPH_DEFAULT_DIRECTORY_SCOPE, roleAssignmentRequest.DirectoryScopeId, "expected role assignment request directoryScopeId to be %s, got %s", GRAPH_DEFAULT_DIRECTORY_SCOPE, roleAssignmentRequest.DirectoryScopeId)
	assert.Equal(t, roleAssignment.RoleDefinitionId, roleAssignmentRequest.RoleDefinitionId, "expected role assignment request roleDefinitionId to be %s, got %s", roleAssignment.RoleDefinitionId, roleAssignmentRequest.RoleDefinitionId)
	assert.Equal(t, expectedDuration, roleAssignmentRequest.ScheduleInfo.Expiration.Duration, "expected role assignment request duration to be %s, got %s", expectedDuration, roleAssignmentRequest.ScheduleInfo.Expiration.Duration)
	assert.Equal(t, StatusProvisioned, requestResponse.Status, "expected role assignment response status to be %s, got %s", StatusProvisioned, requestResponse.Status)
}
