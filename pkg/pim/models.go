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

type PIMResponse struct {
	Name string `json:"name"`
}

type RoleAssignmentSubject struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	DisplayName   string `json:"displayName"`
	PrincipalName string `json:"principalName"`
	Email         string `json:"email"`
}

type RoleResource struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
}

type RoleDefinition struct {
	Id          string        `json:"id"`
	ResourceId  string        `json:"resourceId"`
	Type        string        `json:"type"`
	DisplayName string        `json:"displayName"`
	Resource    *RoleResource `json:"resource"`
}

type RoleAssignment struct {
	Id               string                 `json:"id"`
	ResourceId       string                 `json:"resourceId"`
	RoleDefinitionId string                 `json:"roleDefinitionId"`
	SubjectId        string                 `json:"subjectId"`
	AssignmentState  string                 `json:"assignmentState"`
	Status           string                 `json:"status"`
	Subject          *RoleAssignmentSubject `json:"roleAssignmentSubject"`
	RoleDefinition   *RoleDefinition        `json:"roleDefinition"`
}

type RoleAssignmentResponse struct {
	Value []RoleAssignment `json:"value"`
}

type RoleAssignmentSchedule struct {
	Type          string      `json:"type"`
	StartDateTime interface{} `json:"startDateTime"`
	EndDateTime   interface{} `json:"endDateTime"`
	Duration      string      `json:"duration"`
}

type RoleAssignmentRequest struct {
	RoleDefinitionId               string                  `json:"roleDefinitionId"`
	ResourceId                     string                  `json:"resourceId"`
	SubjectId                      string                  `json:"subjectId"`
	AssignmentState                string                  `json:"assignmentState"`
	Type                           string                  `json:"type"`
	Reason                         string                  `json:"reason"`
	TicketNumber                   string                  `json:"ticketNumber"`
	TicketSystem                   string                  `json:"ticketSystem"`
	Schedule                       *RoleAssignmentSchedule `json:"schedule"`
	LinkedEligibleRoleAssignmentId string                  `json:"linkedEligibleRoleAssignmentId"`
	ScopedResourceId               string                  `json:"scopedResourceId"`
}

type RoleAssignmentRequestStatus struct {
	Status        string              `json:"status"`
	SubStatus     string              `json:"subStatus"`
	StatusDetails []map[string]string `json:"statusDetails"`
}

type RoleAssignmentRequestResponse struct {
	Id                             string                       `json:"id"`
	ResourceId                     string                       `json:"resourceId"`
	RoleDefinitionId               string                       `json:"roleDefinitionId"`
	SubjectId                      string                       `json:"subjectId"`
	ScopedResourceId               string                       `json:"scopedResourceId"`
	LinkedEligibleRoleAssignmentId string                       `json:"linkedEligibleRoleAssignmentId"`
	Type                           string                       `json:"type"`
	AssignmentState                string                       `json:"assignmentState"`
	RequestedDateTime              string                       `json:"requestedDateTime"`
	RoleAssignmentStartDateTime    string                       `json:"roleAssignmentStartDateTime"`
	RoleAssignmentEndDateTime      string                       `json:"roleAssignmentEndDateTime"`
	Reason                         string                       `json:"reason"`
	TicketNumber                   string                       `json:"ticketNumber"`
	TicketSystem                   string                       `json:"ticketSystem"`
	Condition                      string                       `json:"condition"`
	ConditionVersion               string                       `json:"conditionVersion"`
	ConditionDescription           string                       `json:"conditionDescription"`
	Status                         *RoleAssignmentRequestStatus `json:"status"`
	Schedule                       *RoleAssignmentSchedule      `json:"schedule"`
	Metadata                       map[string]interface{}       `json:"metadata"`
}
