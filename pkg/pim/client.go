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
	"log/slog"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netr0m/az-pim-cli/pkg/common"
)

// Azure Client interface
type Client interface {
	GetAccessToken(scope string) string
	GetEligibleResourceAssignments(token string) *ResourceAssignmentResponse
	GetEligibleGovernanceRoleAssignments(roleType string, subjectId string, token string) *GovernanceRoleAssignmentResponse
	ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) bool
	ValidateGovernanceRoleAssignmentRequest(roleType string, roleAssignmentRequest *GovernanceRoleAssignmentRequest, token string) bool
	RequestResourceAssignment(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) *ResourceAssignmentRequestResponse
	RequestGovernanceRoleAssignment(roleType string, governanceRoleAssignmentRequest *GovernanceRoleAssignmentRequest, token string) *GovernanceRoleAssignmentRequestResponse
}

// Azure Client implementation
type AzureClient struct{}

// Implementation of the GetAccessToken call
func (c AzureClient) GetAccessToken(scope string) string {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		_error := common.Error{
			Operation: "GetAccessToken",
			Message:   err.Error(),
			Err:       err,
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}
	tokenOpts := policy.TokenRequestOptions{
		Scopes: []string{
			scope,
		},
	}
	token, err := cred.GetToken(context.Background(), tokenOpts)
	if err != nil {
		_error := common.Error{
			Operation: "GetAccessToken",
			Message:   err.Error(),
			Status:    "401",
			Err:       err,
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	return token.Token
}

func GetAccessToken(scope string, c Client) string {
	return c.GetAccessToken(scope)
}

func GetUserInfo(token string) AzureUserInfo {
	// Decode token
	decoded, err := jwt.ParseWithClaims(token, &AzureUserInfoClaims{}, nil)
	if decoded == nil {
		_error := common.Error{
			Operation: "GetUserInfo",
			Message:   err.Error(),
			Err:       err,
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	// Parse claims
	claims := decoded.Claims.(*AzureUserInfoClaims)

	return *claims.AzureUserInfo
}

func Request(request *PIMRequest, responseModel any) any {
	// Prepare request body
	var req *http.Request
	var err error
	var payload io.Reader
	var _error = common.Error{
		Operation: "Request",
	}

	if request.Payload != nil {
		payload := new(bytes.Buffer)
		json.NewEncoder(payload).Encode(request.Payload) //nolint:errcheck
	} else {
		payload = nil
	}
	req, err = http.NewRequest(request.Method, request.Url, payload)
	if err != nil {
		_error.Message = err.Error()
		_error.Err = err
		_error.Request = req
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		os.Exit(1)
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
		_error.Message = err.Error()
		_error.Status = res.Status
		_error.Err = err
		_error.Request = req
		_error.Response = res
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		os.Exit(1)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		_error.Message = err.Error()
		_error.Status = res.Status
		_error.Err = err
		_error.Request = req
		_error.Response = res
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		os.Exit(1)
	}

	// Handle upstream error responses
	if res.StatusCode >= 400 {
		message := string(body)
		_error.Message = message
		_error.Status = res.Status
		_error.Err = err
		_error.Request = req
		_error.Response = res
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		os.Exit(1)
	}

	err = json.Unmarshal(body, responseModel)
	if err != nil {
		_error.Message = err.Error()
		_error.Status = res.Status
		_error.Err = err
		_error.Request = req
		_error.Response = res
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		os.Exit(1)
	}

	return responseModel
}

func (c AzureClient) GetEligibleResourceAssignments(token string) *ResourceAssignmentResponse {
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

func GetEligibleResourceAssignments(token string, c Client) *ResourceAssignmentResponse {
	return c.GetEligibleResourceAssignments(token)
}

func (c AzureClient) GetEligibleGovernanceRoleAssignments(roleType string, subjectId string, token string) *GovernanceRoleAssignmentResponse {
	if !IsGovernanceRoleType(roleType) {
		_error := common.Error{
			Operation: "GetEligibleGovernanceRoleAssignments",
			Message:   "Invalid role type specified.",
		}
		slog.Error(_error.Error())
		os.Exit(1)
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

func GetEligibleGovernanceRoleAssignments(roleType string, subjectId string, token string, c Client) *GovernanceRoleAssignmentResponse {
	return c.GetEligibleGovernanceRoleAssignments(roleType, subjectId, token)
}

func (c AzureClient) ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) bool {
	var params = map[string]string{
		"api-version": AZ_PIM_API_VERSION,
	}

	resourceAssignmentValidationRequest := resourceAssignmentRequest
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

	return validationResponse.CheckResourceAssignmentResult(resourceAssignmentValidationRequest)
}

func ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string, c Client) bool {
	return c.ValidateResourceAssignmentRequest(scope, resourceAssignmentRequest, token)
}

func (c AzureClient) ValidateGovernanceRoleAssignmentRequest(roleType string, roleAssignmentRequest *GovernanceRoleAssignmentRequest, token string) bool {
	var params = map[string]string{
		"evaluateOnly": "true",
	}

	governanceRoleAssignmentValidationRequest := roleAssignmentRequest

	validationResponse := &GovernanceRoleAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/%s/roleAssignmentRequests", AZ_PIM_GOV_ROLE_BASE_URL, AZ_PIM_GOV_ROLE_BASE_PATH, roleType),
		Token:   token,
		Method:  "POST",
		Params:  params,
		Payload: governanceRoleAssignmentValidationRequest,
	}, validationResponse)

	return validationResponse.CheckGovernanceRoleAssignmentResult(governanceRoleAssignmentValidationRequest)
}

func ValidateGovernanceRoleAssignmentRequest(roleType string, roleAssignmentRequest *GovernanceRoleAssignmentRequest, token string, c Client) bool {
	return c.ValidateGovernanceRoleAssignmentRequest(roleType, roleAssignmentRequest, token)
}

func (c AzureClient) RequestResourceAssignment(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) *ResourceAssignmentRequestResponse {
	var params = map[string]string{
		"api-version": AZ_PIM_API_VERSION,
	}

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

	responseModel.CheckResourceAssignmentResult(resourceAssignmentRequest)

	return responseModel
}

func RequestResourceAssignment(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string, c Client) *ResourceAssignmentRequestResponse {
	return c.RequestResourceAssignment(scope, resourceAssignmentRequest, token)
}

func (c AzureClient) RequestGovernanceRoleAssignment(roleType string, governanceRoleAssignmentRequest *GovernanceRoleAssignmentRequest, token string) *GovernanceRoleAssignmentRequestResponse {
	responseModel := &GovernanceRoleAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/%s/roleAssignmentRequests", AZ_PIM_GOV_ROLE_BASE_URL, AZ_PIM_GOV_ROLE_BASE_PATH, roleType),
		Token:   token,
		Method:  "POST",
		Payload: governanceRoleAssignmentRequest,
	}, responseModel)

	responseModel.CheckGovernanceRoleAssignmentResult(governanceRoleAssignmentRequest)

	return responseModel
}

func RequestGovernanceRoleAssignment(roleType string, governanceRoleAssignmentRequest *GovernanceRoleAssignmentRequest, token string, c Client) *GovernanceRoleAssignmentRequestResponse {
	return c.RequestGovernanceRoleAssignment(roleType, governanceRoleAssignmentRequest, token)
}
