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
	"net/url"
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
	ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) bool
	RequestResourceAssignment(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) *ResourceAssignmentRequestResponse
	GetEligibleGroupAssignments(principalId string, token string) *GraphGroupEligibilityResponse
	GetEligibleRoleAssignments(principalId string, token string) *GraphRoleEligibilityResponse
	RequestGroupAssignment(groupAssignmentRequest *GraphGroupAssignmentRequest, token string) *GraphAssignmentScheduleRequest
	RequestRoleAssignment(roleAssignmentRequest *GraphRoleAssignmentRequest, token string) *GraphAssignmentScheduleRequest
}

// Azure Client implementation
type AzureClient struct {
	ARMBaseURL    string // Azure Resource Manager base URL (Azure resources)
	GraphBaseURL  string // Microsoft Graph base URL (Entra groups and roles)
	GraphScope    string // Microsoft Graph token scope (.default)
	Authority     string // Microsoft Entra authority URL (host + tenant) for sign-in
	ClientID      string // app registration client ID used for Graph access
	TenantID      string // Microsoft Entra tenant ID
	UseDeviceCode bool   // use the device code flow instead of the interactive browser flow
}

// Implementation of the GetAccessToken call
func (c AzureClient) GetAccessToken(scope string) string {
	// Microsoft Graph (Entra groups and roles) is accessed through a custom app
	// registration via an interactive sign-in, since the Azure CLI's built-in
	// client is not authorized for the Graph PIM delegated permissions.
	if c.GraphScope != "" && scope == c.GraphScope {
		return c.getGraphToken(scope)
	}

	// Azure resources (ARM) continue to use the Azure CLI credential
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

	return claims.AzureUserInfo
}

func handleRequestErr(_error *common.Error, err error, req *http.Request) {
	_error.Message = err.Error()
	_error.Err = err
	_error.Request = req
	slog.Error(_error.Error())
	slog.Debug(_error.Debug())
	os.Exit(1)
}

func Request(request *PIMRequest, responseModel any) any {
	// Prepare request body
	var req *http.Request
	var err error
	_error := common.Error{
		Operation: "Request",
	}

	if request.Payload != nil {
		payload := new(bytes.Buffer)
		json.NewEncoder(payload).Encode(request.Payload) //nolint:errcheck
		req, err = http.NewRequest(request.Method, request.Url, payload)
		if err != nil {
			handleRequestErr(&_error, err, req)
		}
	} else {
		req, err = http.NewRequest(request.Method, request.Url, nil)
		if err != nil {
			handleRequestErr(&_error, err, req)
		}
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Token))
	for key, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Prepare request parameters
	query := req.URL.Query()
	for k, v := range request.Params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		defer func() {
			if err := res.Body.Close(); err != nil {
				slog.Error(fmt.Sprintf("Failed to close response body: %v", err))
			}
		}()
		_error.Message = err.Error()
		_error.Status = res.Status
		_error.Err = err
		_error.Request = req
		_error.Response = res
		slog.Error(_error.Error())
		slog.Debug(_error.Debug())
		os.Exit(1)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error(fmt.Sprintf("Failed to close response body: %v", err))
		}
	}()

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
	params := map[string]string{
		"api-version": AZ_PIM_API_VERSION,
		"$filter":     "asTarget()",
	}
	responseModel := &ResourceAssignmentResponse{}
	_ = Request(&PIMRequest{
		Url:    fmt.Sprintf("%s/%s/roleEligibilityScheduleInstances", c.ARMBaseURL, ARM_BASE_PATH),
		Token:  token,
		Method: "GET",
		Params: params,
	}, responseModel)

	return responseModel
}

func GetEligibleResourceAssignments(token string, c Client) *ResourceAssignmentResponse {
	return c.GetEligibleResourceAssignments(token)
}

func (c AzureClient) GetEligibleGroupAssignments(principalId string, token string) *GraphGroupEligibilityResponse {
	params := map[string]string{
		"$filter": fmt.Sprintf("principalId eq '%s'", principalId),
		"$expand": "group",
	}
	responseModel := &GraphGroupEligibilityResponse{}
	_ = Request(&PIMRequest{
		Url:    fmt.Sprintf("%s/%s/%s", c.GraphBaseURL, GRAPH_API_VERSION, GRAPH_GROUP_ELIGIBILITY_PATH),
		Token:  token,
		Method: "GET",
		Params: params,
	}, responseModel)

	return responseModel
}

func GetEligibleGroupAssignments(principalId string, token string, c Client) *GraphGroupEligibilityResponse {
	return c.GetEligibleGroupAssignments(principalId, token)
}

func (c AzureClient) GetEligibleRoleAssignments(principalId string, token string) *GraphRoleEligibilityResponse {
	params := map[string]string{
		"$filter": fmt.Sprintf("principalId eq '%s'", principalId),
		"$expand": "roleDefinition",
	}
	responseModel := &GraphRoleEligibilityResponse{}
	_ = Request(&PIMRequest{
		Url:    fmt.Sprintf("%s/%s/%s", c.GraphBaseURL, GRAPH_API_VERSION, GRAPH_ROLE_ELIGIBILITY_PATH),
		Token:  token,
		Method: "GET",
		Params: params,
	}, responseModel)

	return responseModel
}

func GetEligibleRoleAssignments(principalId string, token string, c Client) *GraphRoleEligibilityResponse {
	return c.GetEligibleRoleAssignments(principalId, token)
}

func (c AzureClient) ValidateResourceAssignmentRequest(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) bool {
	u, err := url.JoinPath(c.ARMBaseURL, scope, ARM_BASE_PATH, "roleAssignmentScheduleRequests", uuid.NewString(), "validate")
	if err != nil {
		_error := common.Error{
			Operation: "ValidateResourceAssignmentRequest",
			Message:   err.Error(),
			Err:       err,
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	params := map[string]string{
		"api-version": AZ_PIM_API_VERSION,
	}

	resourceAssignmentValidationRequest := resourceAssignmentRequest
	resourceAssignmentValidationRequest.Properties.IsValidationOnly = true

	validationResponse := &ResourceAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     u,
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

func (c AzureClient) RequestResourceAssignment(scope string, resourceAssignmentRequest *ResourceAssignmentRequestRequest, token string) *ResourceAssignmentRequestResponse {
	u, err := url.JoinPath(c.ARMBaseURL, scope, ARM_BASE_PATH, "roleAssignmentScheduleRequests", uuid.NewString())
	if err != nil {
		_error := common.Error{
			Operation: "RequestResourceAssignment",
			Message:   err.Error(),
			Err:       err,
		}
		slog.Error(_error.Error())
		os.Exit(1)
	}

	params := map[string]string{
		"api-version": AZ_PIM_API_VERSION,
	}

	responseModel := &ResourceAssignmentRequestResponse{}
	_ = Request(&PIMRequest{
		Url:     u,
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

func (c AzureClient) RequestGroupAssignment(groupAssignmentRequest *GraphGroupAssignmentRequest, token string) *GraphAssignmentScheduleRequest {
	responseModel := &GraphAssignmentScheduleRequest{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/%s", c.GraphBaseURL, GRAPH_API_VERSION, GRAPH_GROUP_REQUEST_PATH),
		Token:   token,
		Method:  "POST",
		Payload: groupAssignmentRequest,
	}, responseModel)

	responseModel.CheckGraphRequestResult()

	return responseModel
}

func RequestGroupAssignment(groupAssignmentRequest *GraphGroupAssignmentRequest, token string, c Client) *GraphAssignmentScheduleRequest {
	return c.RequestGroupAssignment(groupAssignmentRequest, token)
}

func (c AzureClient) RequestRoleAssignment(roleAssignmentRequest *GraphRoleAssignmentRequest, token string) *GraphAssignmentScheduleRequest {
	responseModel := &GraphAssignmentScheduleRequest{}
	_ = Request(&PIMRequest{
		Url:     fmt.Sprintf("%s/%s/%s", c.GraphBaseURL, GRAPH_API_VERSION, GRAPH_ROLE_REQUEST_PATH),
		Token:   token,
		Method:  "POST",
		Payload: roleAssignmentRequest,
	}, responseModel)

	responseModel.CheckGraphRequestResult()

	return responseModel
}

func RequestRoleAssignment(roleAssignmentRequest *GraphRoleAssignmentRequest, token string, c Client) *GraphAssignmentScheduleRequest {
	return c.RequestRoleAssignment(roleAssignmentRequest, token)
}
