/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/

package pim

import (
	"fmt"
)

// Dummy JWT
const TEST_DUMMY_JWT string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI2NjIxYjg4OC02YTc5LTRiNjEtYjlmOS1kMzI1YzUxZWE3OWEiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jYmVkZTYzOC1hM2Q5LTQ1OWYtOGY0ZS0yNGNlZDczYjRlNWUvIiwiaWF0IjoxNzMwMTE0NzE2LCJuYmYiOjE3MzAxMTQ3MTYsImV4cCI6MTczMDExODk2NCwiYWNyIjoiMSIsImFpbyI6IkFiUUFTLzlZQUFBQXVURTdHSVBZMmJrTTI4RjJIUnBtRExiSW5tUDZRTitKY3BEcDBrMnZOQ3J0QlNyZHlwbHBVMlZ0QWpxRUdTTWRrNFRtc2RjMUZWU09HNGtWL0pDaW9ydjl3eDJGT3Zob1M5cVBJSS9FcE8rN3d1dHFMYlFrd3JGd09DU2gvYnU3RXpKV1hKTEJUb1Y3Mk5RS20xc0lPeHNmdm9xQkIvbk5ZRy9SNTFEV1VZRG9vUmdsZ0R0cUduNDNEbCthN2licFNadHRudWM2ejV0VUNCQzFicEczK0lkS2lQOVY4WGhTYThtNHNNL2pWZ2s9IiwiYW1yIjpbInJzYSIsIm1mYSJdLCJhcHBpZCI6ImE0ODQ4Y2I0LTE2YTYtNDRhMi05N2QzLWU1YWUxZmRlMjQ4YSIsImFwcGlkYWNyIjoiMCIsImRldmljZWlkIjoiYmI0OTc2OTgtYWUyNS00YjdiLTkxYjktY2FlNTQxY2M4NTlmIiwiZmFtaWx5X25hbWUiOiJhei1waW0tY2xpIiwiZ2l2ZW5fbmFtZSI6ImF6LXBpbS1jbGktdGVzdCIsImdyb3VwcyI6WyIzNjk5MjNiOC0wZTY2LTRlZWEtODViNS01Y2E5OTk3Y2JkYjciXSwiaWR0eXAiOiJ1c2VyIiwiaXBhZGRyIjoiMTI3LjAuMC4xIiwibmFtZSI6ImF6LXBpbS1jbGkgdGVzdCIsIm9pZCI6IjY1OTVlMTgzLWUzMzMtNDY0Ni1hMWQxLWJjZDA2ZGIzZjIxMiIsIm9ucHJlbV9zaWQiOiJTLTEtNS0yMS0xMzM3MDAxMzMzLTEzMzczNjgyNS0xMzM3NDU1NDMtNTkxMzM3IiwicHVpZCI6IjEwMDEzMzczQkMyQkFFRDciLCJyaCI6IjAuQVM5QVBQYnR5OW1qbjBXUFRpVE8xenRPWHFjel9BRzdlRDBOcEdkMmpqTnVpUTR2QU5zLiIsInNjcCI6InVzZXJfaW1wZXJzb25hdGlvbiIsInN1YiI6IlBkWE5UNnhaUUtqb1pSeldKZWRJeEZXUTZzR2g2cU9JWmNjcjJYMC1MckEiLCJ0aWQiOiJkOTc2ZDI1Yy1kOWE1LTRhMTEtYWFhYy04Nzc0YjY3NjYyOGEiLCJ1bmlxdWVfbmFtZSI6ImF6LXBpbS1jbGktdGVzdEBuZXRyMG0uZ2l0aHViLmNvbSIsInVwbiI6ImF6LXBpbS1jbGktdGVzdEBuZXRyMG0uZ2l0aHViLmNvbSIsInV0aSI6IkN0SnhueHF6SGttR0NoZUNGT0ZaQkIiLCJ2ZXIiOiIxLjAiLCJ3aWRzIjpbImFiY2QxMzM3LTEzMzctMTMzNy0xMzM3LTc2YjE5NGU4NTUwOSJdLCJ4bXNfaWRyZWwiOiIxIDE2IiwieG1zX3RkYnIiOiJFVSJ9.0vjNDIYJirD2FUF8oGIVNw_Q4_VC432qRj3EBqmPeqU"

// Dummy principal/subject ID
const TEST_DUMMY_PRINCIPAL_ID string = "6595e183-e333-4646-a1d1-bcd06db3f212"

// Dummy principal/subject name
const TEST_DUMMY_PRINCIPAL_NAME string = "az-pim-cli test"

// Dummy principal/subject name
const TEST_DUMMY_PRINCIPAL_EMAIL string = "az-pim-cli-test@netr0m.github.com"

// Dummy role names
const (
	TEST_DUMMY_ROLE_1_NAME = "Role 1"
	TEST_DUMMY_ROLE_2_NAME = "Role 2"
)

// Dummy subscription data
const (
	TEST_DUMMY_SUBSCRIPTION_1_NAME = "Subscription 1"
	TEST_DUMMY_SUBSCRIPTION_1_ID   = "6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83"
	TEST_DUMMY_SUB_1_ROLE_1_ID     = "e7721be1-2fd0-4dc7-adff-4cc0de716f6f"
	TEST_DUMMY_SUB_1_ROLE_2_ID     = "eccdef45-ccd4-426c-937c-d1a3a43cad8e"
	TEST_DUMMY_SUBSCRIPTION_2_NAME = "Subscription 2"
	TEST_DUMMY_SUBSCRIPTION_2_ID   = "1472abc0-1996-46b0-a068-f52770a7ddc5"
	TEST_DUMMY_SUB_2_ROLE_1_ID     = "fc0b56ad-4a56-482a-bd50-eb20d84e1d61"
	TEST_DUMMY_SUBSCRIPTION_3_NAME = "Azure Resource 1"
	TEST_DUMMY_SUBSCRIPTION_3_ID   = "2e7c848b-f6d1-44f1-92b8-b1034a1d038d"
	TEST_DUMMY_SUB_3_ROLE_1_ID     = "eed43a65-3dc1-4040-a094-c0145429c729"
)

// Dummy AAD Group data
const (
	TEST_DUMMY_GROUP_1_NAME    = "Group 1"
	TEST_DUMMY_GROUP_1_ID      = "50c01cd1-0bd1-42bd-8820-01e4342bb90c"
	TEST_DUMMY_GRP_1_ROLE_1_ID = "29b079fb-7f51-42da-bb91-a8f451ff0e0f"
	TEST_DUMMY_GRP_1_ROLE_2_ID = "7d26ca21-db44-48ec-8346-0ff53a483ccf"
	TEST_DUMMY_GROUP_2_NAME    = "Group 2"
	TEST_DUMMY_GROUP_2_ID      = "f1177909-f650-4acc-94e4-8de3691e91fe"
	TEST_DUMMY_GRP_2_ROLE_1_ID = "edb0a3b5-133e-4e2d-aaf9-566b033c615a"
)

var EligibleResourceAssignmentsDummyData *ResourceAssignmentResponse = &ResourceAssignmentResponse{
	Value: []ResourceAssignment{
		{
			Id:   "3638336a-34b7-4275-a34e-096d37cac30e",
			Name: fmt.Sprintf("%s %s", TEST_DUMMY_SUBSCRIPTION_1_NAME, TEST_DUMMY_ROLE_1_NAME),
			Type: "Subscription",
			Properties: &ResourceProperties{
				ExpandedProperties: &ResourceExpandedProperties{
					Scope: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUBSCRIPTION_1_ID,
						DisplayName: TEST_DUMMY_SUBSCRIPTION_1_NAME,
					},
					RoleDefinition: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUB_1_ROLE_1_ID,
						DisplayName: TEST_DUMMY_ROLE_1_NAME,
					},
					Principal: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_PRINCIPAL_ID,
						DisplayName: TEST_DUMMY_PRINCIPAL_NAME,
						Email:       TEST_DUMMY_PRINCIPAL_EMAIL,
					},
				},
			},
		},
		{
			Id:   "cb758738-4954-4608-ad23-45087d510e6f",
			Name: fmt.Sprintf("%s %s", TEST_DUMMY_SUBSCRIPTION_1_NAME, TEST_DUMMY_ROLE_2_NAME),
			Type: "Subscription",
			Properties: &ResourceProperties{
				ExpandedProperties: &ResourceExpandedProperties{
					Scope: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUBSCRIPTION_1_ID,
						DisplayName: TEST_DUMMY_SUBSCRIPTION_1_NAME,
					},
					RoleDefinition: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUB_1_ROLE_2_ID,
						DisplayName: TEST_DUMMY_ROLE_2_NAME,
					},
					Principal: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_PRINCIPAL_ID,
						DisplayName: TEST_DUMMY_PRINCIPAL_NAME,
						Email:       TEST_DUMMY_PRINCIPAL_EMAIL,
					},
				},
			},
		},
		{
			Id:   "c973c4e3-1afc-4879-82a1-eae6447414b8",
			Name: fmt.Sprintf("%s %s", TEST_DUMMY_SUBSCRIPTION_2_NAME, TEST_DUMMY_ROLE_1_NAME),
			Type: "Subscription",
			Properties: &ResourceProperties{
				ExpandedProperties: &ResourceExpandedProperties{
					Scope: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUBSCRIPTION_2_ID,
						DisplayName: TEST_DUMMY_SUBSCRIPTION_2_NAME,
					},
					RoleDefinition: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUB_2_ROLE_1_ID,
						DisplayName: TEST_DUMMY_ROLE_1_NAME,
					},
					Principal: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_PRINCIPAL_ID,
						DisplayName: TEST_DUMMY_PRINCIPAL_NAME,
						Email:       TEST_DUMMY_PRINCIPAL_EMAIL,
					},
				},
			},
		},
		{
			Id:   "33fde24b-903b-4960-8ff2-8d611ec50ac3",
			Name: fmt.Sprintf("%s %s", TEST_DUMMY_SUBSCRIPTION_3_NAME, TEST_DUMMY_ROLE_1_NAME),
			Type: "Subscription",
			Properties: &ResourceProperties{
				ExpandedProperties: &ResourceExpandedProperties{
					Scope: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUBSCRIPTION_3_ID,
						DisplayName: TEST_DUMMY_SUBSCRIPTION_3_NAME,
					},
					RoleDefinition: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_SUB_3_ROLE_1_ID,
						DisplayName: TEST_DUMMY_ROLE_1_NAME,
					},
					Principal: &ResourceExpandedProperty{
						Id:          TEST_DUMMY_PRINCIPAL_ID,
						DisplayName: TEST_DUMMY_PRINCIPAL_NAME,
						Email:       TEST_DUMMY_PRINCIPAL_EMAIL,
					},
				},
			},
		},
	},
}

// Dummy Entra role data
const (
	TEST_DUMMY_ROLE_1_DEFINITION_ID = "729827e3-9c14-49f7-bb1b-9608f156bbb8"
	TEST_DUMMY_ROLE_2_DEFINITION_ID = "fe930be7-5e62-47db-91af-98c3a49a38b1"
)

var EligibleGroupAssignmentsDummyData *GraphGroupEligibilityResponse = &GraphGroupEligibilityResponse{
	Value: []GraphGroupEligibilityInstance{
		{
			Id:          "38622a10-cc6c-48cf-a516-9ebbffcc5aae",
			PrincipalId: TEST_DUMMY_PRINCIPAL_ID,
			GroupId:     TEST_DUMMY_GROUP_1_ID,
			AccessId:    "member",
			Group: &GraphGroup{
				Id:          TEST_DUMMY_GROUP_1_ID,
				DisplayName: TEST_DUMMY_GROUP_1_NAME,
			},
		},
		{
			Id:          "2f6203dd-ddc5-49c3-9b1f-fa6b9f6e1d4d",
			PrincipalId: TEST_DUMMY_PRINCIPAL_ID,
			GroupId:     TEST_DUMMY_GROUP_1_ID,
			AccessId:    "owner",
			Group: &GraphGroup{
				Id:          TEST_DUMMY_GROUP_1_ID,
				DisplayName: TEST_DUMMY_GROUP_1_NAME,
			},
		},
		{
			Id:          "533faa7b-5a9c-483f-b854-c01e380b4ae7",
			PrincipalId: TEST_DUMMY_PRINCIPAL_ID,
			GroupId:     TEST_DUMMY_GROUP_2_ID,
			AccessId:    "member",
			Group: &GraphGroup{
				Id:          TEST_DUMMY_GROUP_2_ID,
				DisplayName: TEST_DUMMY_GROUP_2_NAME,
			},
		},
	},
}

var EligibleRoleAssignmentsDummyData *GraphRoleEligibilityResponse = &GraphRoleEligibilityResponse{
	Value: []GraphRoleEligibilityInstance{
		{
			Id:               "a3f9b3b1-2b6c-4d7e-8f1a-1c2d3e4f5a6b",
			PrincipalId:      TEST_DUMMY_PRINCIPAL_ID,
			RoleDefinitionId: TEST_DUMMY_ROLE_1_DEFINITION_ID,
			DirectoryScopeId: "/",
			RoleDefinition: &GraphRoleDefinition{
				Id:          TEST_DUMMY_ROLE_1_DEFINITION_ID,
				DisplayName: TEST_DUMMY_ROLE_1_NAME,
			},
		},
		{
			Id:               "b4a0c4c2-3c7d-4e8f-9a2b-2d3e4f5a6b7c",
			PrincipalId:      TEST_DUMMY_PRINCIPAL_ID,
			RoleDefinitionId: TEST_DUMMY_ROLE_2_DEFINITION_ID,
			DirectoryScopeId: "/",
			RoleDefinition: &GraphRoleDefinition{
				Id:          TEST_DUMMY_ROLE_2_DEFINITION_ID,
				DisplayName: TEST_DUMMY_ROLE_2_NAME,
			},
		},
	},
}
