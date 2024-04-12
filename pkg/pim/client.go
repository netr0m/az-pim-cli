/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v5"
)

func GetPIMAccessTokenAzureCLI() string {
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		log.Fatalln(err)
	}
	tokenOpts := policy.TokenRequestOptions{
		Scopes: []string{
			AZ_PIM_SCOPE,
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
	url := fmt.Sprintf("%s/%s", AZ_PIM_BASE_URL, request.Path)

	// Prepare request body
	var req *http.Request
	var err error
	if request.Payload != nil {
		payload := new(bytes.Buffer)
		json.NewEncoder(payload).Encode(request.Payload)
		req, err = http.NewRequest(request.Method, url, payload)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
	} else {
		// Prepare the request
		req, err = http.NewRequest(request.Method, url, nil)
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
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
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

func GetEligibleRoleAssignments(token string) *RoleAssignmentResponse {
	var params = map[string]string{
		"api-version": AZ_PIM_API_VERSION,
		"$filter":     "asTarget()",
	}
	responseModel := &RoleAssignmentResponse{}
	_ = Request(&PIMRequest{
		Path:   fmt.Sprintf("%s/roleEligibilityScheduleInstances", AZ_PIM_BASE_PATH),
		Token:  token,
		Method: "GET",
		Params: params,
	}, responseModel)

	return responseModel
}

func RequestRoleAssignment(subjectId string, subscriptionId string, roleDefinitionId string, roleAssignmentId string, duration int, token string, resourceType string, reason string) *RoleAssignmentRequestResponse {
	responseModel := &RoleAssignmentRequestResponse{}
	roleAssignmentSchedule := &RoleAssignmentSchedule{
		Type:          "Once",
		StartDateTime: nil,
		EndDateTime:   nil,
		Duration:      fmt.Sprintf("PT%dM", duration),
	}
	roleAssignmentRequest := &RoleAssignmentRequest{
		RoleDefinitionId:               roleDefinitionId,
		ResourceId:                     subscriptionId,
		SubjectId:                      subjectId,
		AssignmentState:                "Active",
		Type:                           "UserAdd",
		Reason:                         reason,
		TicketNumber:                   "",
		TicketSystem:                   "az-pim-cli",
		Schedule:                       roleAssignmentSchedule,
		LinkedEligibleRoleAssignmentId: roleAssignmentId,
		ScopedResourceId:               "",
	}

	_ = Request(&PIMRequest{
		Path:    fmt.Sprintf("%s/roleAssignmentRequests", resourceType),
		Token:   token,
		Method:  "POST",
		Payload: roleAssignmentRequest,
	}, responseModel)

	return responseModel
}
