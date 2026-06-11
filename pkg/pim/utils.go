/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package pim

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/netr0m/az-pim-cli/pkg/common"
)

func IsResourceAssignmentRequestFailed(requestResponse *ResourceAssignmentRequestResponse) bool {
	switch requestResponse.Properties.Status {
	case StatusAdminDenied, StatusCanceled, StatusDenied, StatusFailed, StatusFailedAsResourceIsLocked, StatusInvalid, StatusRevoked, StatusTimedOut:
		return true
	}
	return false
}

func IsGraphRequestFailed(requestResponse *GraphAssignmentScheduleRequest) bool {
	switch requestResponse.Status {
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

func IsGraphRequestPending(requestResponse *GraphAssignmentScheduleRequest) bool {
	switch requestResponse.Status {
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

func IsGraphRequestOK(requestResponse *GraphAssignmentScheduleRequest) bool {
	switch requestResponse.Status {
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

func (response *GraphAssignmentScheduleRequest) CheckGraphRequestResult() bool {
	if IsGraphRequestFailed(response) {
		_error := common.Error{
			Operation: "CheckGraphRequestResult",
			Message:   "The role assignment request failed",
			Status:    response.Status,
			Response:  response,
		}
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		return false
	}
	if IsGraphRequestOK(response) {
		slog.Info("The role assignment request was successful", "status", response.Status)
		return true
	}
	if IsGraphRequestPending(response) {
		slog.Warn("The role assignment request is pending", "status", response.Status)
		return true
	}

	return false
}

func parseDate(dateStr string) (time.Time, error) {
	dateLayout := "02/01/2006" // DD/MM/YYYY
	d, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return d, err
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local), nil
}

func parseTime(timeStr string) (time.Time, error) {
	timeLayout := "15:04" // HH:MM
	return time.Parse(timeLayout, timeStr)
}

func parseDateTime(dateStr string, timeStr string) (string, *common.Error) {
	_error := &common.Error{
		Operation: "parseDateTime",
	}
	var d time.Time
	if dateStr == "" {
		// Get the current date, remove time info
		now := time.Now().Local()
		d = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		_d, err := parseDate(dateStr)
		if err != nil {
			_error.Message = "Unable to parse the date"
			_error.Err = err
			return "", _error
		}
		d = _d
	}
	if timeStr != "" {
		t, err := parseTime(timeStr)
		if err != nil {
			_error.Message = "Unable to parse the time"
			_error.Err = err
			return "", _error
		}
		d = d.Add(time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute()))
	}

	// format is the same as 'time.RFC3339', except for the timezone part.
	// must be customized due to inconsistencies in behavior when TZ environment variable is missing
	formatted := d.Format("2006-01-02T15:04:05-07:00")
	return formatted, nil
}

func CreateResourceAssignmentScheduleInfo(duration int, startDate string, startTime string) *ScheduleInfo {
	var scheduleStart interface{}
	if (startDate != "") || (startTime != "") {
		startDateTime, err := parseDateTime(startDate, startTime)
		if err != nil {
			slog.Error(err.Error())
			slog.Debug(err.Debug())
			os.Exit(1)
		}
		scheduleStart = startDateTime
	}

	return &ScheduleInfo{
		StartDateTime: scheduleStart,
		Expiration: &ScheduleInfoExpiration{
			Type:     "AfterDuration",
			Duration: fmt.Sprintf("PT%dM", duration),
		},
	}
}

func CreateResourceAssignmentRequest(subjectId string, resourceAssignment *ResourceAssignment, duration int, startDate string, startTime string, reason string, ticketSystem string, ticketNumber string) (string, *ResourceAssignmentRequestRequest) {
	return CreateResourceAssignmentRequestWithScope(subjectId, resourceAssignment, "", duration, startDate, startTime, reason, ticketSystem, ticketNumber)
}

func CreateResourceAssignmentRequestWithScope(subjectId string, resourceAssignment *ResourceAssignment, scope string, duration int, startDate string, startTime string, reason string, ticketSystem string, ticketNumber string) (string, *ResourceAssignmentRequestRequest) {
	scheduleInfo := CreateResourceAssignmentScheduleInfo(duration, startDate, startTime)
	resourceAssignmentRequest := &ResourceAssignmentRequestRequest{
		Properties: ResourceAssignmentRequestProperties{
			PrincipalId:                     subjectId,
			RoleDefinitionId:                resourceAssignment.Properties.ExpandedProperties.RoleDefinition.Id,
			RequestType:                     "SelfActivate",
			LinkedRoleEligibilityScheduleId: resourceAssignment.Properties.RoleEligibilityScheduleId,
			Justification:                   reason,
			ScheduleInfo:                    scheduleInfo,
			TicketInfo:                      &TicketInfo{TicketNumber: ticketNumber, TicketSystem: ticketSystem},
			IsValidationOnly:                false,
			IsActivativation:                true,
		},
	}
	if scope == "" {
		scope = resourceAssignment.Properties.ExpandedProperties.Scope.Id
	}

	return scope, resourceAssignmentRequest
}

func CreateGraphScheduleInfo(duration int, startDate string, startTime string) *GraphScheduleInfo {
	var startDateTime *string
	if (startDate != "") || (startTime != "") {
		s, err := parseDateTime(startDate, startTime)
		if err != nil {
			slog.Error(err.Error())
			slog.Debug(err.Debug())
			os.Exit(1)
		}
		startDateTime = &s
	}

	return &GraphScheduleInfo{
		StartDateTime: startDateTime,
		Expiration: &GraphScheduleInfoExpiration{
			Type:     GRAPH_EXPIRATION_AFTER_DURATION,
			Duration: fmt.Sprintf("PT%dM", duration),
		},
	}
}

func CreateGraphGroupAssignmentRequest(principalId string, instance *GraphGroupEligibilityInstance, duration int, startDate string, startTime string, reason string, ticketSystem string, ticketNumber string) *GraphGroupAssignmentRequest {
	return &GraphGroupAssignmentRequest{
		Action:        GRAPH_ACTION_SELF_ACTIVATE,
		AccessId:      instance.AccessId,
		PrincipalId:   principalId,
		GroupId:       instance.GroupId,
		Justification: reason,
		ScheduleInfo:  CreateGraphScheduleInfo(duration, startDate, startTime),
		TicketInfo:    &GraphTicketInfo{TicketNumber: ticketNumber, TicketSystem: ticketSystem},
	}
}

func CreateGraphRoleAssignmentRequest(principalId string, instance *GraphRoleEligibilityInstance, duration int, startDate string, startTime string, reason string, ticketSystem string, ticketNumber string) *GraphRoleAssignmentRequest {
	directoryScopeId := instance.DirectoryScopeId
	if directoryScopeId == "" {
		directoryScopeId = GRAPH_DEFAULT_DIRECTORY_SCOPE
	}

	return &GraphRoleAssignmentRequest{
		Action:           GRAPH_ACTION_SELF_ACTIVATE,
		PrincipalId:      principalId,
		RoleDefinitionId: instance.RoleDefinitionId,
		DirectoryScopeId: directoryScopeId,
		Justification:    reason,
		ScheduleInfo:     CreateGraphScheduleInfo(duration, startDate, startTime),
		TicketInfo:       &GraphTicketInfo{TicketNumber: ticketNumber, TicketSystem: ticketSystem},
	}
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

func (instance *GraphGroupEligibilityInstance) Debug() string {
	var debugLines []string

	debugLines = append(debugLines, fmt.Sprintf("ID: %s", instance.Id))
	debugLines = append(debugLines, fmt.Sprintf("\tGroupId: %s", instance.GroupId))
	debugLines = append(debugLines, fmt.Sprintf("\tAccessId: %s", instance.AccessId))
	debugLines = append(debugLines, fmt.Sprintf("\tPrincipalId: %s", instance.PrincipalId))
	if instance.Group != nil {
		debugLines = append(debugLines, fmt.Sprintf("\tGroup: %s", instance.Group.DisplayName))
	}

	return strings.Join(debugLines, "\n")
}

func (instance *GraphRoleEligibilityInstance) Debug() string {
	var debugLines []string

	debugLines = append(debugLines, fmt.Sprintf("ID: %s", instance.Id))
	debugLines = append(debugLines, fmt.Sprintf("\tRoleDefinitionId: %s", instance.RoleDefinitionId))
	debugLines = append(debugLines, fmt.Sprintf("\tDirectoryScopeId: %s", instance.DirectoryScopeId))
	debugLines = append(debugLines, fmt.Sprintf("\tPrincipalId: %s", instance.PrincipalId))
	if instance.RoleDefinition != nil {
		debugLines = append(debugLines, fmt.Sprintf("\tRoleDefinition: %s", instance.RoleDefinition.DisplayName))
	}

	return strings.Join(debugLines, "\n")
}
