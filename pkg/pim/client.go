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
	"log/slog"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netr0m/az-pim-cli/pkg/common"
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
		error := common.Error{
			Operation: "GetPIMAccessTokenAzureCLI",
			Message:   err.Error(),
			Status:    "401",
			Err:       err,
		}
		slog.Error(error.Error())
		os.Exit(1)
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

func GetEligibleGovernanceRoleAssignments(roleType string, subjectId string, token string) *GovernanceRoleAssignmentResponse {
	if !IsGovernanceRoleType(roleType) {
		log.Fatalln("GetEligibleGovernanceRoleAssignments: Invalid role type specified.")
	}
	var params = map[string]string{
		"$expand": "linkedEligibleRoleAssignment,subject,scopedResource,roleDefinition($expand=resource)",
		"$filter": fmt.Sprintf("(subject/id eq '%s') and (assignmentState eq 'Eligible')", subjectId),
	}
	responseModel := &GovernanceRoleAssignmentResponse{}
	_ = Request(&PIMRequest{
		Url:    fmt.Sprintf("%s/%s/%s/roleAssignments", AZ_PIM_GOV_ROLE_BASE_URL, AZ_PIM_GOV_ROLE_BASE_PATH, roleType),
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

func ValidateGovernanceRoleAssignmentRequest(roleType string, roleAssignmentRequest GovernanceRoleAssignmentRequest, token string) bool {
	if !IsGovernanceRoleType(roleType) {
		log.Fatalln("ValidateGovernanceRoleAssignmentRequest: Invalid role type specified.")
	}
	var params = map[string]string{
		"evaluateOnly": "true",
	}

	governanceRoleAssignmentValidationRequest := roleAssignmentRequest
	governanceRoleAssignmentValidationRequest.Reason = "Evaluate Only"
	governanceRoleAssignmentValidationRequest.TicketNumber = "Evaluate Only"
	governanceRoleAssignmentValidationRequest.TicketSystem = "Evaluate Only"

	validationResponse := &GovernanceRoleAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/%s/roleAssignmentRequests", AZ_PIM_GOV_ROLE_BASE_URL, AZ_PIM_GOV_ROLE_BASE_PATH, roleType),
		Token:   token,
		Method:  "POST",
		Params:  params,
		Payload: governanceRoleAssignmentValidationRequest,
	}, validationResponse)

	if IsGovernanceRoleAssignmentRequestFailed(validationResponse) {
		log.Printf("ERROR: The role assignment validation failed with status '%s', '%s'", validationResponse.Status.Status, validationResponse.Status.SubStatus)
		log.Fatalln(validationResponse)
		return false
	}
	if IsGovernanceRoleAssignmentRequestOK(validationResponse) {
		return true
	}
	if IsGovernanceRoleAssignmentRequestPending(validationResponse) {
		log.Printf("WARNING: The role assignment request is pending with status '%s', '%s'", validationResponse.Status.Status, validationResponse.Status.SubStatus)
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

func RequestGovernanceRoleAssignment(subjectId string, roleType string, governanceRoleAssignment *GovernanceRoleAssignment, duration int, reason string, ticketSystem string, ticketNumber string, token string) *GovernanceRoleAssignmentRequestResponse {
	if !IsGovernanceRoleType(roleType) {
		log.Fatalln("RequestGovernanceRoleAssignment: Invalid role type specified.")
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

	ValidateGovernanceRoleAssignmentRequest(roleType, *governanceRoleAssignmentRequest, token)

	responseModel := &GovernanceRoleAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/%s/roleAssignmentRequests", AZ_PIM_GOV_ROLE_BASE_URL, AZ_PIM_GOV_ROLE_BASE_PATH, roleType),
		Token:   token,
		Method:  "POST",
		Payload: governanceRoleAssignmentRequest,
	}, responseModel)

	return responseModel
}
