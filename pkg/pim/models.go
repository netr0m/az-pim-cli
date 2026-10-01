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

type GovernanceRoleAssignmentSubject struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	DisplayName   string `json:"displayName"`
	PrincipalName string `json:"principalName"`
	Email         string `json:"email"`
}

type GovernanceRoleResource struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
}

type GovernanceRoleDefinition struct {
	Id          string                  `json:"id"`
	ResourceId  string                  `json:"resourceId"`
	Type        string                  `json:"type"`
	DisplayName string                  `json:"displayName"`
	Resource    *GovernanceRoleResource `json:"resource"`
}

type GovernanceRoleAssignment struct {
	Id               string                           `json:"id"`
	ResourceId       string                           `json:"resourceId"`
	RoleDefinitionId string                           `json:"roleDefinitionId"`
	AccessId         string                           `json:"accessId"`
	SubjectId        string                           `json:"subjectId"`
	AssignmentState  string                           `json:"assignmentState"`
	Status           string                           `json:"status"`
	Subject          *GovernanceRoleAssignmentSubject `json:"subject"`
	RoleDefinition   *GovernanceRoleDefinition        `json:"roleDefinition"`
}

type GovernanceRoleAssignmentResponse struct {
	Value []GovernanceRoleAssignment `json:"value"`
}

type graphRoleDefinition struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type graphPrincipal struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type graphRoleEligibilityScheduleInstance struct {
	Id                        string               `json:"id"`
	PrincipalId               string               `json:"principalId"`
	RoleDefinitionId          string               `json:"roleDefinitionId"`
	DirectoryScopeId          string               `json:"directoryScopeId"`
	RoleEligibilityScheduleId string               `json:"roleEligibilityScheduleId"`
	RoleDefinition            *graphRoleDefinition `json:"roleDefinition"`
	Principal                 *graphPrincipal      `json:"principal"`
}

type graphRoleEligibilityScheduleInstanceResponse struct {
	Value []graphRoleEligibilityScheduleInstance `json:"value"`
}

type graphGroup struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type graphGroupEligibilityScheduleInstance struct {
	Id                    string          `json:"id"`
	PrincipalId           string          `json:"principalId"`
	GroupId               string          `json:"groupId"`
	AccessId              string          `json:"accessId"`
	EligibilityScheduleId string          `json:"eligibilityScheduleId"`
	Group                 *graphGroup     `json:"group"`
	Principal             *graphPrincipal `json:"principal"`
}

type graphGroupEligibilityScheduleInstanceResponse struct {
	Value []graphGroupEligibilityScheduleInstance `json:"value"`
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
	EndDateTime   interface{}             `json:"endDateTime,omitempty"`
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

type GovernanceRoleAssignmentRequest struct {
	Action           string        `json:"action"`
	PrincipalId      string        `json:"principalId"`
	RoleDefinitionId string        `json:"roleDefinitionId,omitempty"`
	DirectoryScopeId string        `json:"directoryScopeId,omitempty"`
	GroupId          string        `json:"groupId,omitempty"`
	AccessId         string        `json:"accessId,omitempty"`
	Justification    string        `json:"justification"`
	ScheduleInfo     *ScheduleInfo `json:"scheduleInfo"`
	TicketInfo       *TicketInfo   `json:"ticketInfo"`
	IsValidationOnly bool          `json:"isValidationOnly"`
}

type GovernanceRoleAssignmentRequestResponse struct {
	Id               string        `json:"id"`
	Status           string        `json:"status"`
	PrincipalId      string        `json:"principalId"`
	RoleDefinitionId string        `json:"roleDefinitionId,omitempty"`
	DirectoryScopeId string        `json:"directoryScopeId,omitempty"`
	GroupId          string        `json:"groupId,omitempty"`
	AccessId         string        `json:"accessId,omitempty"`
	IsValidationOnly bool          `json:"isValidationOnly"`
	TargetScheduleId string        `json:"targetScheduleId"`
	Justification    string        `json:"justification"`
	ScheduleInfo     *ScheduleInfo `json:"scheduleInfo"`
	TicketInfo       *TicketInfo   `json:"ticketInfo"`
}
