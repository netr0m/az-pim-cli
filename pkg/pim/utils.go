/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package pim

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/netr0m/az-pim-cli/pkg/common"
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

func (response *ResourceAssignmentRequestResponse) CheckResourceAssignmentResult(request *ResourceAssignmentRequestRequest) bool {
	if IsResourceAssignmentRequestFailed(response) {
		_error := common.Error{
			Operation: "CheckResourceAssignmentResult",
			Message:   "The role assignment validation failed",
			Status:    response.Properties.Status,
			Request:   request,
			Response:  response,
		}
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		return false
	}
	if IsResourceAssignmentRequestOK(response) {
		slog.Info("The role assignment request was successful", "status", response.Properties.Status)
		return true
	}
	if IsResourceAssignmentRequestPending(response) {
		slog.Warn("The role assignment request is pending", "status", response.Properties.Status)
		return true
	}

	return false
}

func (response *GovernanceRoleAssignmentRequestResponse) CheckGovernanceRoleAssignmentResult(request *GovernanceRoleAssignmentRequest) bool {
	if IsGovernanceRoleAssignmentRequestFailed(response) {
		_error := common.Error{
			Operation: "CheckGovernanceRoleAssignmentResult",
			Message:   "The role assignment validation failed",
			Status:    response.Status.Status,
			Request:   request,
			Response:  response,
		}
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		return false
	}
	if IsGovernanceRoleAssignmentRequestOK(response) {
		slog.Info("The role assignment request was successfully validated", "status", response.Status.Status, "subStatus", response.Status.SubStatus)
		return true
	}
	if IsGovernanceRoleAssignmentRequestPending(response) {
		slog.Warn("The role assignment request is pending", "status", response.Status.Status, "subStatus", response.Status.SubStatus)
		return true
	}

	return false
}

func CreateResourceAssignmentRequest(subjectId string, resourceAssignment *ResourceAssignment, duration int, reason string, ticketSystem string, ticketNumber string) (string, *ResourceAssignmentRequestRequest) {
	resourceAssignmentRequest := &ResourceAssignmentRequestRequest{
		Properties: ResourceAssignmentRequestProperties{
			PrincipalId:                     subjectId,
			RoleDefinitionId:                resourceAssignment.Properties.ExpandedProperties.RoleDefinition.Id,
			RequestType:                     "SelfActivate",
			LinkedRoleEligibilityScheduleId: resourceAssignment.Properties.RoleEligibilityScheduleId,
			Justification:                   reason,
			ScheduleInfo: &ScheduleInfo{
				StartDateTime: nil,
				Expiration: &ScheduleInfoExpiration{
					Type:     "AfterDuration",
					Duration: fmt.Sprintf("PT%dM", duration),
				},
			},
			TicketInfo:       &TicketInfo{TicketNumber: ticketNumber, TicketSystem: ticketSystem},
			IsValidationOnly: false,
			IsActivativation: true,
		},
	}
	scope := resourceAssignment.Properties.ExpandedProperties.Scope.Id[1:]

	return scope, resourceAssignmentRequest
}

func CreateGovernanceRoleAssignmentRequest(subjectId string, roleType string, governanceRoleAssignment *GovernanceRoleAssignment, duration int, reason string, ticketSystem string, ticketNumber string) (string, *GovernanceRoleAssignmentRequest) {
	if !IsGovernanceRoleType(roleType) {
		_error := common.Error{
			Operation: "CreateGovernanceRoleAssignmentRequest",
			Message:   "Invalid role type specified.",
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	governanceRoleAssignmentRequest := &GovernanceRoleAssignmentRequest{
		RoleDefinitionId: governanceRoleAssignment.RoleDefinitionId,
		ResourceId:       governanceRoleAssignment.ResourceId,
		SubjectId:        subjectId,
		AssignmentState:  "Active",
		Type:             "UserAdd",
		Reason:           reason,
		TicketNumber:     ticketNumber,
		TicketSystem:     ticketSystem,
		Schedule: &GovernanceRoleAssignmentSchedule{
			Type:          "Once",
			StartDateTime: nil,
			EndDateTime:   nil,
			Duration:      fmt.Sprintf("PT%dM", duration),
		},
		LinkedEligibleRoleAssignmentId: governanceRoleAssignment.Id,
		ScopedResourceId:               "",
	}

	return roleType, governanceRoleAssignmentRequest
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
