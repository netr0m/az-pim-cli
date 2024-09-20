/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetPIMAccessTokenAzureCLI(scope string) string {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		log.Fatalln(err)
	}
	tokenOpts := policy.TokenRequestOptions{
		Scopes: []string{
			scope,
		},
	}
	token, err := cred.GetToken(context.Background(), tokenOpts)
	if err != nil {
		log.Fatalln(err)
	}

	return token.Token
}

func GetUserInfo(token string) AzureUserInfo {
	// Decode token
	decoded, err := jwt.ParseWithClaims(token, &AzureUserInfoClaims{}, nil)
	if decoded == nil {
		log.Fatalln(err)
	}

	// Parse claims
	claims := decoded.Claims.(*AzureUserInfoClaims)

	return *claims.AzureUserInfo
}

func Request(request *PIMRequest, responseModel any) any {
	// Prepare request body
	var req *http.Request
	var err error
	if request.Payload != nil {
		payload := new(bytes.Buffer)
		json.NewEncoder(payload).Encode(request.Payload) //nolint:errcheck
		req, err = http.NewRequest(request.Method, request.Url, payload)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
	} else {
		// Prepare the request
		req, err = http.NewRequest(request.Method, request.Url, nil)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
	}
	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Token))

	// Prepare request parameters
	query := req.URL.Query()
	for k, v := range request.Params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		defer res.Body.Close()
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Handle upstream error responses
	if res.StatusCode >= 400 {
		message := string(body)
		log.Fatalf("The upstream API responded with status %s: %s\n", res.Status, message)
	}

	err = json.Unmarshal(body, responseModel)
	if err != nil {
		log.Fatalln(err)
	}

	return responseModel
}

func GetEligibleResourceAssignments(token string) *ResourceAssignmentResponse {
	var params = map[string]string{
		"api-version": AZ_PIM_API_VERSION,
		"$filter":     "asTarget()",
	}
	responseModel := &ResourceAssignmentResponse{}
	_ = Request(&PIMRequest{
		Url:    fmt.Sprintf("%s/%s/roleEligibilityScheduleInstances", AZ_PIM_BASE_URL, AZ_PIM_BASE_PATH),
		Token:  token,
		Method: "GET",
		Params: params,
	}, responseModel)

	return responseModel
}

func GetEligibleGroupAssignments(token string, subjectId string) *GroupAssignmentResponse {
	var params = map[string]string{
		"$expand": "linkedEligibleRoleAssignment,subject,scopedResource,roleDefinition($expand=resource)",
		"$filter": fmt.Sprintf("(subject/id eq '%s') and (assignmentState eq 'Eligible')", subjectId),
	}
	responseModel := &GroupAssignmentResponse{}
	_ = Request(&PIMRequest{
		Url:    fmt.Sprintf("%s/%s/aadGroups/roleAssignments", AZ_PIM_GROUP_BASE_URL, AZ_PIM_GROUP_BASE_PATH),
		Token:  token,
		Method: "GET",
		Params: params,
	}, responseModel)

	return responseModel
}

func ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest ResourceAssignmentRequestRequest, token string) bool {
	var params = map[string]string{
		"api-version": AZ_PIM_API_VERSION,
	}

	resourceAssignmentValidationRequest := resourceAssignmentRequest
	resourceAssignmentValidationRequest.Properties.Justification = "validation only call"
	resourceAssignmentValidationRequest.Properties.TicketInfo.TicketNumber = "Evaluate Only"
	resourceAssignmentValidationRequest.Properties.TicketInfo.TicketSystem = "Evaluate Only"
	resourceAssignmentValidationRequest.Properties.IsValidationOnly = true

	validationResponse := &ResourceAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url: fmt.Sprintf(
			"%s/%s/%s/roleAssignmentScheduleRequests/%s/validate",
			AZ_PIM_BASE_URL,
			scope,
			AZ_PIM_BASE_PATH,
			uuid.NewString(),
		),
		Token:   token,
		Method:  "POST",
		Params:  params,
		Payload: resourceAssignmentValidationRequest,
	}, validationResponse)

	if IsResourceAssignmentRequestFailed(validationResponse) {
		log.Printf("ERROR: The role assignment validation failed with status '%s'", validationResponse.Properties.Status)
		log.Fatalln(validationResponse)
		return false
	}
	if IsResourceAssignmentRequestOK(validationResponse) {
		return true
	}
	if IsResourceAssignmentRequestPending(validationResponse) {
		log.Printf("WARNING: The role assignment request is pending with status '%s'", validationResponse.Properties.Status)
		return true
	}

	return false
}

func ValidateGroupAssignmentRequest(groupAssignmentRequest GroupAssignmentRequest, token string) bool {
	var params = map[string]string{
		"evaluateOnly": "true",
	}

	groupAssignmentValidationRequest := groupAssignmentRequest
	groupAssignmentValidationRequest.Reason = "Evaluate Only"
	groupAssignmentValidationRequest.TicketNumber = "Evaluate Only"
	groupAssignmentValidationRequest.TicketSystem = "Evaluate Only"

	validationResponse := &GroupAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/aadGroups/roleAssignmentRequests", AZ_PIM_GROUP_BASE_URL, AZ_PIM_GROUP_BASE_PATH),
		Token:   token,
		Method:  "POST",
		Params:  params,
		Payload: groupAssignmentValidationRequest,
	}, validationResponse)

	if IsGroupAssignmentRequestFailed(validationResponse) {
		log.Printf("ERROR: The group assignment validation failed with status '%s', '%s'", validationResponse.Status.Status, validationResponse.Status.SubStatus)
		log.Fatalln(validationResponse)
		return false
	}
	if IsGroupAssignmentRequestOK(validationResponse) {
		return true
	}
	if IsGroupAssignmentRequestPending(validationResponse) {
		log.Printf("WARNING: The group assignment request is pending with status '%s', '%s'", validationResponse.Status.Status, validationResponse.Status.SubStatus)
		return true
	}

	return false
}

func RequestResourceAssignment(subjectId string, resourceAssignment *ResourceAssignment, duration int, reason string, ticketSystem string, ticketNumber string, token string) *ResourceAssignmentRequestResponse {
	var params = map[string]string{
		"api-version": AZ_PIM_API_VERSION,
	}

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

	ValidateResourceAssignmentRequest(scope, *resourceAssignmentRequest, token)

	responseModel := &ResourceAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url: fmt.Sprintf(
			"%s/%s/%s/roleAssignmentScheduleRequests/%s",
			AZ_PIM_BASE_URL,
			scope,
			AZ_PIM_BASE_PATH,
			uuid.NewString(),
		),
		Token:   token,
		Method:  "PUT",
		Params:  params,
		Payload: resourceAssignmentRequest,
	}, responseModel)

	return responseModel
}

func RequestGroupAssignment(subjectId string, groupAssignment *GroupAssignment, duration int, reason string, ticketSystem string, ticketNumber string, token string) *GroupAssignmentRequestResponse {
	groupAssignmentRequest := &GroupAssignmentRequest{
		RoleDefinitionId: groupAssignment.RoleDefinitionId,
		ResourceId:       groupAssignment.ResourceId,
		SubjectId:        subjectId,
		AssignmentState:  "Active",
		Type:             "UserAdd",
		Reason:           reason,
		TicketNumber:     ticketNumber,
		TicketSystem:     ticketSystem,
		Schedule: &GroupAssignmentSchedule{
			Type:          "Once",
			StartDateTime: nil,
			EndDateTime:   nil,
			Duration:      fmt.Sprintf("PT%dM", duration),
		},
		LinkedEligibleRoleAssignmentId: groupAssignment.Id,
		ScopedResourceId:               "",
	}

	ValidateGroupAssignmentRequest(*groupAssignmentRequest, token)

	responseModel := &GroupAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/aadGroups/roleAssignmentRequests", AZ_PIM_GROUP_BASE_URL, AZ_PIM_GROUP_BASE_PATH),
		Token:   token,
		Method:  "POST",
		Payload: groupAssignmentRequest,
	}, responseModel)

	return responseModel
}
