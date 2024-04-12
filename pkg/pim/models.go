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
	*jwt.MapClaims
	*AzureUserInfo
	extra map[string]interface{}
}

type PIMRequest struct {
	Path    string
	Token   string
	Method  string
	Headers map[string][]string
	Payload interface{}
	Params  map[string]string
}

type RoleExpandedProperty struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Email       string `json:"email"`
}

type RoleExpandedProperties struct {
	Principal      *RoleExpandedProperty `json:"principal"`
	RoleDefinition *RoleExpandedProperty `json:"roleDefinition"`
	Scope          *RoleExpandedProperty `json:"scope"`
}

type RoleProperties struct {
	RoleEligibilityScheduleId string                  `json:"roleEligibilityScheduleId"`
	Scope                     string                  `json:"scope"`
	RoleDefinitionId          string                  `json:"roleDefinitionId"`
	PrincipalId               string                  `json:"principalId"`
	PrincipalType             string                  `json:"principalType"`
	Status                    string                  `json:"status"`
	StartDateTime             string                  `json:"startDateTime"`
	EndDateTime               string                  `json:"endDateTime"`
	MemberType                string                  `json:"memberType"`
	CreatedOn                 string                  `json:"createdOn"`
	Condition                 string                  `json:"condition"`
	ConditionVersion          string                  `json:"conditionVersion"`
	ExpandedProperties        *RoleExpandedProperties `json:"expandedProperties"`
}

type RoleAssignment struct {
	Properties *RoleProperties `json:"properties"`
	Name       string          `json:"name"`
	Id         string          `json:"id"`
	Type       string          `json:"type"`
}

type RoleAssignmentResponse struct {
	Value []RoleAssignment `json:"value"`
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

type RoleAssignmentValidationProperties struct {
	LinkedRoleEligibilityScheduleId string                  `json:"linkedRoleEligibilityScheduleId"`
	TargetRoleAssignmentScheduleId  string                  `json:"targetRoleAssignmentScheduleId"`
	Scope                           string                  `json:"scope"`
	RoleDefinitionId                string                  `json:"roleDefinitionId"`
	PrincipalId                     string                  `json:"principalId"`
	PrincipalType                   string                  `json:"principalType"`
	RequestType                     string                  `json:"requestType"`
	Status                          string                  `json:"status"`
	ScheduleInfo                    *ScheduleInfo           `json:"scheduleInfo"`
	TicketInfo                      *TicketInfo             `json:"ticketInfo"`
	Justification                   string                  `json:"justification"`
	RequestorId                     string                  `json:"requestorId"`
	CreatedOn                       string                  `json:"createdOn"`
	ExpandedProperties              *RoleExpandedProperties `json:"expandedProperties"`
}

type RoleAssignmentRequestResponse struct {
	Properties *RoleAssignmentValidationProperties `json:"properties"`
	Name       string                              `json:"name"`
	Id         string                              `json:"id"`
	Type       string                              `json:"type"`
}

type RoleAssignmentRequestProperties struct {
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

type RoleAssignmentRequestRequest struct {
	Properties RoleAssignmentRequestProperties `json:"Properties"`
}
