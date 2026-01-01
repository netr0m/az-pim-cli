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

func (m *mockClient) GetEligibleGovernanceRoleAssignments(roleType string, subjectId string, token string) *GovernanceRoleAssignmentResponse {
	args := m.Called(roleType, subjectId, token)
	return args.Get(0).(*GovernanceRoleAssignmentResponse)
}

func TestGetEligibleGovernanceRoleAssignmentsAADGroup(t *testing.T) {
	m := newMockClient()

	m.On("GetEligibleGovernanceRoleAssignments", ROLE_TYPE_AAD_GROUPS, TEST_DUMMY_PRINCIPAL_ID, TEST_DUMMY_JWT).Return(EligibleGovernanceRoleAssignmentsDummyData)

	eligibleGovernanceRoleAssignments := GetEligibleGovernanceRoleAssignments(ROLE_TYPE_AAD_GROUPS, TEST_DUMMY_PRINCIPAL_ID, TEST_DUMMY_JWT, m)

	if len(eligibleGovernanceRoleAssignments.Value) != 3 {
		t.Errorf("expected 3 eligible governance role assignments, got %v", len(eligibleGovernanceRoleAssignments.Value))
	}
	for _, governanceRole := range eligibleGovernanceRoleAssignments.Value {
		if governanceRole.SubjectId != TEST_DUMMY_PRINCIPAL_ID {
			t.Errorf("expected governance role SubjectId to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, governanceRole.SubjectId)
		}
	}
	// Check group name
	_groupName := eligibleGovernanceRoleAssignments.Value[1].RoleDefinition.Resource.DisplayName
	if _groupName != TEST_DUMMY_GROUP_1_NAME {
		t.Errorf("expected governance role RoleDefinition.Resource.DisplayName to be %s, got %s", TEST_DUMMY_GROUP_1_NAME, _groupName)
	}
	// Check role name
	_roleName := eligibleGovernanceRoleAssignments.Value[2].RoleDefinition.DisplayName
	if _roleName != TEST_DUMMY_ROLE_1_NAME {
		t.Errorf("expected governance role RoleDefinition.DisplayName to be %s, got %s", TEST_DUMMY_ROLE_1_NAME, _roleName)
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

func (m *mockClient) ValidateGovernanceRoleAssignmentRequest(roleType string, roleAssignmentRequest *GovernanceRoleAssignmentRequest, token string) bool {
	args := m.Called(roleType, roleAssignmentRequest, token)
	return args.Bool(0)
}

func TestValidateGovernanceRoleAssignmentRequest(t *testing.T) {
	m := newMockClient()

	governanceRoleAssignment := &EligibleGovernanceRoleAssignmentsDummyData.Value[0]
	roleType, governanceRoleAssignmentRequest := CreateGovernanceRoleAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, ROLE_TYPE_AAD_GROUPS, governanceRoleAssignment, 30, "", "", "test", "Test", "1337")

	m.On("ValidateGovernanceRoleAssignmentRequest", roleType, governanceRoleAssignmentRequest, TEST_DUMMY_JWT).Return(true)

	isValid := ValidateGovernanceRoleAssignmentRequest(roleType, governanceRoleAssignmentRequest, TEST_DUMMY_JWT, m)

	if !isValid {
		t.Errorf("expected governance role assignment request validation to be successful, got %v", isValid)
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

func (m *mockClient) RequestGovernanceRoleAssignment(roleType string, governanceRoleAssignmentRequest *GovernanceRoleAssignmentRequest, token string) *GovernanceRoleAssignmentRequestResponse {
	args := m.Called(roleType, governanceRoleAssignmentRequest, token)
	return args.Get(0).(*GovernanceRoleAssignmentRequestResponse)
}

func TestRequestGovernanceRoleAssignmentAADGroup(t *testing.T) {
	m := newMockClient()

	governanceRoleAssignment := &EligibleGovernanceRoleAssignmentsDummyData.Value[0]
	roleType, governanceRoleAssignmentRequest := CreateGovernanceRoleAssignmentRequest(TEST_DUMMY_PRINCIPAL_ID, ROLE_TYPE_AAD_GROUPS, governanceRoleAssignment, DEFAULT_DURATION_MINUTES, "", "", DEFAULT_REASON, "Test", "1337")
	governanceRoleAssignmentRequestResponse := &GovernanceRoleAssignmentRequestResponse{
		Id:               governanceRoleAssignment.Id,
		ResourceId:       governanceRoleAssignmentRequest.ResourceId,
		RoleDefinitionId: governanceRoleAssignmentRequest.RoleDefinitionId,
		SubjectId:        governanceRoleAssignment.SubjectId,
		AssignmentState:  governanceRoleAssignmentRequest.AssignmentState,
		Status: &GovernanceRoleAssignmentRequestStatus{
			Status:    "Active",
			SubStatus: "Active",
		},
		TicketSystem:                   "Test",
		TicketNumber:                   "1337",
		Reason:                         DEFAULT_REASON,
		Schedule:                       governanceRoleAssignmentRequest.Schedule,
		LinkedEligibleRoleAssignmentId: governanceRoleAssignmentRequest.LinkedEligibleRoleAssignmentId,
		ScopedResourceId:               governanceRoleAssignmentRequest.ScopedResourceId,
	}

	m.On("RequestGovernanceRoleAssignment", ROLE_TYPE_AAD_GROUPS, governanceRoleAssignmentRequest, TEST_DUMMY_JWT).Return(governanceRoleAssignmentRequestResponse)

	requestResponse := RequestGovernanceRoleAssignment(roleType, governanceRoleAssignmentRequest, TEST_DUMMY_JWT, m)
	expectedDuration := fmt.Sprintf("PT%dM", DEFAULT_DURATION_MINUTES)

	assert.Equal(t, requestResponse.Reason, DEFAULT_REASON, "expected governance role assignment request reason to be %s, got %s", DEFAULT_REASON, requestResponse.Reason)
	assert.Equal(t, requestResponse.SubjectId, TEST_DUMMY_PRINCIPAL_ID, "expected governance role assignment request subject ID to be %s, got %s", TEST_DUMMY_PRINCIPAL_ID, requestResponse.SubjectId)
	assert.Equal(t, requestResponse.Status.Status, "Active", "expected governance role assignment request status to be %s, got %s", "Active", requestResponse.Status.Status)
	assert.Equal(t, requestResponse.Schedule.Duration, expectedDuration, "expected governance role assignment request expiration duration to be %s, got %s", expectedDuration, requestResponse.Schedule.Duration)
}
