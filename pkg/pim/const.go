/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

/*
* US Gov Cloud: https://learn.microsoft.com/en-us/azure/azure-government/compare-azure-government-global-azure#guidance-for-developers
* China: https://learn.microsoft.com/en-us/azure/china/concepts-service-availability#azure-in-china-rest-endpoints
 */

// Base URL for the Azure Resource Manager (ARM) API (Azure Resources)
const (
	ARM_GLOBAL_BASE_URL = "https://management.azure.com"
	ARM_USGOV_BASE_URL  = "https://management.usgovcloudapi.net"
	ARM_CN_BASE_URL     = "https://management.chinacloudapi.cn"
)

// Base path for the Azure Resource Manager PIM API
const ARM_BASE_PATH string = "providers/Microsoft.Authorization"

// Base URLs for Microsoft Graph (Entra Groups and Entra Roles) per Azure environment.
// See https://learn.microsoft.com/en-us/graph/deployments
const (
	GRAPH_GLOBAL_BASE_URL = "https://graph.microsoft.com"
	GRAPH_USGOV_BASE_URL  = "https://graph.microsoft.us"
	GRAPH_CN_BASE_URL     = "https://microsoftgraph.chinacloudapi.cn"
)

// Microsoft Entra authority hosts per Azure environment, used for MSAL device code sign-in.
// See https://learn.microsoft.com/en-us/entra/identity-platform/authentication-national-cloud
const (
	GRAPH_GLOBAL_AUTHORITY_HOST = "login.microsoftonline.com"
	GRAPH_USGOV_AUTHORITY_HOST  = "login.microsoftonline.us"
	GRAPH_CN_AUTHORITY_HOST     = "login.partner.microsoftonline.cn"
)

// API version for Microsoft Graph
const GRAPH_API_VERSION string = "v1.0"

// Default scope requested for Microsoft Graph tokens.
const GRAPH_DEFAULT_SCOPE string = "https://graph.microsoft.com/.default"

// Default reason for role activation
const DEFAULT_REASON string = "config"

// Default duration for role activation
const DEFAULT_DURATION_MINUTES int = 480

// API version for the "role eligibility schedule instances" (i.e. eligible azure resource role assignments)
const AZ_PIM_API_VERSION string = "2020-10-01"

// Role types
const (
	ROLE_TYPE_AAD_GROUPS  = "aadGroups"
	ROLE_TYPE_ENTRA_ROLES = "aadroles"
)

// Base URLs for different Azure environments
var ARM_BASE_URLS = map[string]string{
	"global": ARM_GLOBAL_BASE_URL,
	"usgov":  ARM_USGOV_BASE_URL,
	"china":  ARM_CN_BASE_URL,
}

// Microsoft Graph base URLs for different Azure environments
var GRAPH_BASE_URLS = map[string]string{
	"global": GRAPH_GLOBAL_BASE_URL,
	"usgov":  GRAPH_USGOV_BASE_URL,
	"china":  GRAPH_CN_BASE_URL,
}

// Microsoft Entra authority hosts for different Azure environments
var GRAPH_AUTHORITY_HOSTS = map[string]string{
	"global": GRAPH_GLOBAL_AUTHORITY_HOST,
	"usgov":  GRAPH_USGOV_AUTHORITY_HOST,
	"china":  GRAPH_CN_AUTHORITY_HOST,
}
