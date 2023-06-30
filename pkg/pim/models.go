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
