/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package pim

import (
	"fmt"
	"strings"
)

func IsResourceAssignmentRequestFailed(requestResponse *ResourceAssignmentRequestResponse) bool {
	switch requestResponse.Properties.Status {
	case StatusAdminDenied, StatusCanceled, StatusDenied, StatusFailed, StatusFailedAsResourceIsLocked, StatusInvalid, StatusRevoked, StatusTimedOut:
		return true
	}
	return false
}

func IsGovernanceRoleAssignmentRequestFailed(requestResponse *GovernanceRoleAssignmentRequestResponse) bool {
	switch requestResponse.Status.SubStatus {
	case StatusAdminDenied, StatusCanceled, StatusDenied, StatusFailed, StatusFailedAsResourceIsLocked, StatusInvalid, StatusRevoked, StatusTimedOut:
		return true
	}
	return false
}

func IsResourceAssignmentRequestPending(requestResponse *ResourceAssignmentRequestResponse) bool {
	switch requestResponse.Properties.Status {
	case StatusPendingAdminDecision, StatusPendingApproval, StatusPendingApprovalProvisioning, StatusPendingEvaluation, StatusPendingExternalProvisioning, StatusPendingProvisioning, StatusPendingRevocation, StatusPendingScheduleCreation:
		return true
	}
	return false
}

func IsGovernanceRoleAssignmentRequestPending(requestResponse *GovernanceRoleAssignmentRequestResponse) bool {
	switch requestResponse.Status.SubStatus {
	case StatusPendingAdminDecision, StatusPendingApproval, StatusPendingApprovalProvisioning, StatusPendingEvaluation, StatusPendingExternalProvisioning, StatusPendingProvisioning, StatusPendingRevocation, StatusPendingScheduleCreation:
		return true
	}
	return false
}

func IsResourceAssignmentRequestOK(requestResponse *ResourceAssignmentRequestResponse) bool {
	switch requestResponse.Properties.Status {
	case StatusAccepted, StatusAdminApproved, StatusGranted, StatusProvisioned, StatusProvisioningStarted, StatusScheduleCreated:
		return true
	}
	return false
}

func IsGovernanceRoleAssignmentRequestOK(requestResponse *GovernanceRoleAssignmentRequestResponse) bool {
	switch requestResponse.Status.SubStatus {
	case StatusAccepted, StatusAdminApproved, StatusGranted, StatusProvisioned, StatusProvisioningStarted, StatusScheduleCreated:
		return true
	}
	return false
}

func IsGovernanceRoleType(roleType string) bool {
	switch roleType {
	case ROLE_TYPE_AAD_GROUPS, ROLE_TYPE_ENTRA_ROLES:
		return true
	}
	return false
}

func (resourceAssignment *ResourceAssignment) Debug() string {
	var debugLines []string

	debugLines = append(debugLines, fmt.Sprintf("ID: %s", resourceAssignment.Id))
	if resourceAssignment.Properties != nil {
		if resourceAssignment.Properties.ExpandedProperties != nil {
			debugLines = append(debugLines, fmt.Sprintf("\tScopeID: %s", resourceAssignment.Properties.ExpandedProperties.Scope.Id))
			if resourceAssignment.Properties.ExpandedProperties.Principal != nil {
				debugLines = append(debugLines, fmt.Sprintf("\tPrincipal: %s", resourceAssignment.Properties.ExpandedProperties.Principal.DisplayName))
			}
			if resourceAssignment.Properties.ExpandedProperties.RoleDefinition != nil {
				debugLines = append(debugLines, fmt.Sprintf("\tRoleDefinition: %s", resourceAssignment.Properties.ExpandedProperties.RoleDefinition.DisplayName))
			}
		}
		debugLines = append(debugLines, fmt.Sprintf("\tRoleDefinitionId: %s", resourceAssignment.Properties.RoleDefinitionId))
		debugLines = append(debugLines, fmt.Sprintf("\tPrincipalID: %s", resourceAssignment.Properties.PrincipalId))
		debugLines = append(debugLines, fmt.Sprintf("\tStatus: %s", resourceAssignment.Properties.Status))
	}

	return strings.Join(debugLines, "\n")
}

func (roleAssignment *GovernanceRoleAssignment) Debug() string {
	var debugLines []string

	debugLines = append(debugLines, fmt.Sprintf("ID: %s", roleAssignment.Id))
	debugLines = append(debugLines, fmt.Sprintf("\tResourceID: %s", roleAssignment.ResourceId))
	debugLines = append(debugLines, fmt.Sprintf("\tRoleDefinitionId: %s", roleAssignment.RoleDefinitionId))
	debugLines = append(debugLines, fmt.Sprintf("\tSubjectId: %s", roleAssignment.SubjectId))
	debugLines = append(debugLines, fmt.Sprintf("\tAssignmentState: %s", roleAssignment.AssignmentState))
	debugLines = append(debugLines, fmt.Sprintf("\tStatus: %s", roleAssignment.Status))
	if roleAssignment.Subject != nil {
		debugLines = append(debugLines, fmt.Sprintf("\tSubject: %s", roleAssignment.Subject.DisplayName))
	}
	if roleAssignment.RoleDefinition != nil {
		debugLines = append(debugLines, fmt.Sprintf("\tRoleDefinition: %s", roleAssignment.RoleDefinition.DisplayName))
	}

	return strings.Join(debugLines, "\n")
}
