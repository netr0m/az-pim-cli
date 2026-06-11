/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

import "github.com/golang-jwt/jwt/v5"

type AzureUserInfo struct {
	ObjectId string `json:"oid"`
	Email    string `json:"unique_name"`
}

type AzureUserInfoClaims struct {
	jwt.MapClaims
	AzureUserInfo
}

type PIMRequest struct {
	Url     string
	Token   string
	Method  string
	Headers map[string][]string
	Payload interface{}
	Params  map[string]string
}

type ResourceExpandedProperty struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Email       string `json:"email"`
}

type ResourceExpandedProperties struct {
	Principal      *ResourceExpandedProperty `json:"principal"`
	RoleDefinition *ResourceExpandedProperty `json:"roleDefinition"`
	Scope          *ResourceExpandedProperty `json:"scope"`
}

type ResourceProperties struct {
	RoleEligibilityScheduleId string                      `json:"roleEligibilityScheduleId"`
	Scope                     string                      `json:"scope"`
	RoleDefinitionId          string                      `json:"roleDefinitionId"`
	PrincipalId               string                      `json:"principalId"`
	PrincipalType             string                      `json:"principalType"`
	Status                    string                      `json:"status"`
	StartDateTime             string                      `json:"startDateTime"`
	EndDateTime               string                      `json:"endDateTime"`
	MemberType                string                      `json:"memberType"`
	CreatedOn                 string                      `json:"createdOn"`
	Condition                 string                      `json:"condition"`
	ConditionVersion          string                      `json:"conditionVersion"`
	ExpandedProperties        *ResourceExpandedProperties `json:"expandedProperties"`
}

type ResourceAssignment struct {
	Properties *ResourceProperties `json:"properties"`
	Name       string              `json:"name"`
	Id         string              `json:"id"`
	Type       string              `json:"type"`
}

type ResourceAssignmentResponse struct {
	Value []ResourceAssignment `json:"value"`
}

type TicketInfo struct {
	TicketNumber string `json:"ticketNumber"`
	TicketSystem string `json:"ticketSystem"`
}

type ScheduleInfoExpiration struct {
	Type     string `json:"type"`
	Duration string `json:"duration"`
}

type ScheduleInfo struct {
	StartDateTime interface{}             `json:"startDateTime"`
	Expiration    *ScheduleInfoExpiration `json:"expiration"`
	EndDateTime   interface{}             `json:"endDateTime"`
}

const (
	StatusAccepted                    string = "Accepted"
	StatusAdminApproved               string = "AdminApproved"
	StatusAdminDenied                 string = "AdminDenied"
	StatusCanceled                    string = "Canceled"
	StatusDenied                      string = "Denied"
	StatusFailed                      string = "Failed"
	StatusFailedAsResourceIsLocked    string = "FailedAsResourceIsLocked"
	StatusGranted                     string = "Granted"
	StatusInvalid                     string = "Invalid"
	StatusPendingAdminDecision        string = "PendingAdminDecision"
	StatusPendingApproval             string = "PendingApproval"
	StatusPendingApprovalProvisioning string = "PendingApprovalProvisioning"
	StatusPendingEvaluation           string = "PendingEvaluation"
	StatusPendingExternalProvisioning string = "PendingExternalProvisioning"
	StatusPendingProvisioning         string = "PendingProvisioning"
	StatusPendingRevocation           string = "PendingRevocation"
	StatusPendingScheduleCreation     string = "PendingScheduleCreation"
	StatusProvisioned                 string = "Provisioned"
	StatusProvisioningStarted         string = "ProvisioningStarted"
	StatusRevoked                     string = "Revoked"
	StatusScheduleCreated             string = "ScheduleCreated"
	StatusTimedOut                    string = "TimedOut"
)

type ResourceAssignmentValidationProperties struct {
	LinkedRoleEligibilityScheduleId string                      `json:"linkedRoleEligibilityScheduleId"`
	TargetRoleAssignmentScheduleId  string                      `json:"targetRoleAssignmentScheduleId"`
	Scope                           string                      `json:"scope"`
	RoleDefinitionId                string                      `json:"roleDefinitionId"`
	PrincipalId                     string                      `json:"principalId"`
	PrincipalType                   string                      `json:"principalType"`
	RequestType                     string                      `json:"requestType"`
	Status                          string                      `json:"status"`
	ScheduleInfo                    *ScheduleInfo               `json:"scheduleInfo"`
	TicketInfo                      *TicketInfo                 `json:"ticketInfo"`
	Justification                   string                      `json:"justification"`
	RequestorId                     string                      `json:"requestorId"`
	CreatedOn                       string                      `json:"createdOn"`
	ExpandedProperties              *ResourceExpandedProperties `json:"expandedProperties"`
}

type ResourceAssignmentRequestResponse struct {
	Properties *ResourceAssignmentValidationProperties `json:"properties"`
	Name       string                                  `json:"name"`
	Id         string                                  `json:"id"`
	Type       string                                  `json:"type"`
}

type ResourceAssignmentRequestProperties struct {
	PrincipalId                     string        `json:"PrincipalId"`
	RoleDefinitionId                string        `json:"RoleDefinitionId"`
	RequestType                     string        `json:"RequestType"`
	LinkedRoleEligibilityScheduleId string        `json:"LinkedRoleEligibilityScheduleId"`
	Justification                   string        `json:"Justification"`
	ScheduleInfo                    *ScheduleInfo `json:"ScheduleInfo"`
	TicketInfo                      *TicketInfo   `json:"TicketInfo"`
	IsValidationOnly                bool          `json:"IsValidationOnly"`
	IsActivativation                bool          `json:"IsActivativation"` // yes, this typo is in the API
}

type ResourceAssignmentRequestRequest struct {
	Properties ResourceAssignmentRequestProperties `json:"Properties"`
}

// == Microsoft Graph PIM models (Entra groups and roles) ==

// GraphScheduleInfoExpiration describes how an activation expires.
type GraphScheduleInfoExpiration struct {
	Type     string `json:"type"`     // e.g. "afterDuration"
	Duration string `json:"duration"` // ISO-8601 duration, e.g. "PT480M"
}

// GraphScheduleInfo describes when an activation starts and how long it lasts.
// StartDateTime is omitted when not explicitly requested, letting Graph default
// to the time the request is processed.
type GraphScheduleInfo struct {
	StartDateTime *string                      `json:"startDateTime,omitempty"`
	Expiration    *GraphScheduleInfoExpiration `json:"expiration"`
}

// GraphTicketInfo carries optional ticketing metadata for an activation.
type GraphTicketInfo struct {
	TicketNumber string `json:"ticketNumber,omitempty"`
	TicketSystem string `json:"ticketSystem,omitempty"`
}

// GraphGroup is the expanded group object on a group eligibility instance.
type GraphGroup struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// GraphGroupEligibilityInstance is an eligible PIM-for-Groups membership/ownership.
type GraphGroupEligibilityInstance struct {
	Id          string      `json:"id"`
	PrincipalId string      `json:"principalId"`
	GroupId     string      `json:"groupId"`
	AccessId    string      `json:"accessId"` // "member" | "owner"
	Group       *GraphGroup `json:"group"`    // populated via $expand=group
}

type GraphGroupEligibilityResponse struct {
	Value []GraphGroupEligibilityInstance `json:"value"`
}

// GraphRoleDefinition is the expanded role definition on a role eligibility instance.
type GraphRoleDefinition struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// GraphRoleEligibilityInstance is an eligible PIM-for-Entra-roles assignment.
type GraphRoleEligibilityInstance struct {
	Id               string               `json:"id"`
	PrincipalId      string               `json:"principalId"`
	RoleDefinitionId string               `json:"roleDefinitionId"`
	DirectoryScopeId string               `json:"directoryScopeId"`
	RoleDefinition   *GraphRoleDefinition `json:"roleDefinition"` // populated via $expand=roleDefinition
}

type GraphRoleEligibilityResponse struct {
	Value []GraphRoleEligibilityInstance `json:"value"`
}

// GraphGroupAssignmentRequest is the body of a PIM-for-Groups selfActivate request.
type GraphGroupAssignmentRequest struct {
	Action        string             `json:"action"` // "selfActivate"
	AccessId      string             `json:"accessId"`
	PrincipalId   string             `json:"principalId"`
	GroupId       string             `json:"groupId"`
	Justification string             `json:"justification"`
	ScheduleInfo  *GraphScheduleInfo `json:"scheduleInfo"`
	TicketInfo    *GraphTicketInfo   `json:"ticketInfo,omitempty"`
}

// GraphRoleAssignmentRequest is the body of a PIM-for-Entra-roles selfActivate request.
type GraphRoleAssignmentRequest struct {
	Action           string             `json:"action"` // "selfActivate"
	PrincipalId      string             `json:"principalId"`
	RoleDefinitionId string             `json:"roleDefinitionId"`
	DirectoryScopeId string             `json:"directoryScopeId"`
	Justification    string             `json:"justification"`
	ScheduleInfo     *GraphScheduleInfo `json:"scheduleInfo"`
	TicketInfo       *GraphTicketInfo   `json:"ticketInfo,omitempty"`
}

// GraphAssignmentScheduleRequest is the response returned when creating an
// assignment schedule request for either groups or Entra roles.
type GraphAssignmentScheduleRequest struct {
	Id              string `json:"id"`
	Status          string `json:"status"`
	Action          string `json:"action"`
	CreatedDateTime string `json:"createdDateTime"`
}
