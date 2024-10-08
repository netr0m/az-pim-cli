/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package pim

func IsResourceAssignmentRequestFailed(requestResponse *ResourceAssignmentRequestResponse) bool {
	switch requestResponse.Properties.Status {
	case StatusAdminDenied, StatusCanceled, StatusDenied, StatusFailed, StatusFailedAsResourceIsLocked, StatusInvalid, StatusRevoked, StatusTimedOut:
		return true
	}
	return false
}

func IsGroupAssignmentRequestFailed(requestResponse *GroupAssignmentRequestResponse) bool {
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

func IsGroupAssignmentRequestPending(requestResponse *GroupAssignmentRequestResponse) bool {
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

func IsGroupAssignmentRequestOK(requestResponse *GroupAssignmentRequestResponse) bool {
	switch requestResponse.Status.SubStatus {
	case StatusAccepted, StatusAdminApproved, StatusGranted, StatusProvisioned, StatusProvisioningStarted, StatusScheduleCreated:
		return true
	}
	return false
}
