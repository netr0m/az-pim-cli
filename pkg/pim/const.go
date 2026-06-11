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

// Default reason for role activation
const DEFAULT_REASON string = "config"

// Default duration for role activation
const DEFAULT_DURATION_MINUTES int = 480

// API version for the "role eligibility schedule instances" (i.e. eligible azure resource role assignments)
const AZ_PIM_API_VERSION string = "2020-10-01"

// Role types (used as a discriminator between the group and Entra role flows)
const (
	ROLE_TYPE_AAD_GROUPS  = "aadGroups"
	ROLE_TYPE_ENTRA_ROLES = "aadroles"
)

// Base URL for the Microsoft Graph API (Entra Groups and Entra Roles via PIM)
const (
	GRAPH_GLOBAL_BASE_URL = "https://graph.microsoft.com"
	GRAPH_USGOV_BASE_URL  = "https://graph.microsoft.us"
	GRAPH_CN_BASE_URL     = "https://microsoftgraph.chinacloudapi.cn"
)

// Microsoft Graph API version
const GRAPH_API_VERSION string = "v1.0"

// Microsoft Graph PIM API paths
const (
	// PIM for Groups
	GRAPH_GROUP_ELIGIBILITY_PATH = "identityGovernance/privilegedAccess/group/eligibilityScheduleInstances"
	GRAPH_GROUP_REQUEST_PATH     = "identityGovernance/privilegedAccess/group/assignmentScheduleRequests"
	// PIM for Entra (directory) roles
	GRAPH_ROLE_ELIGIBILITY_PATH = "roleManagement/directory/roleEligibilityScheduleInstances"
	GRAPH_ROLE_REQUEST_PATH     = "roleManagement/directory/roleAssignmentScheduleRequests"
)

// Microsoft Graph PIM request action and expiration type
const (
	GRAPH_ACTION_SELF_ACTIVATE      = "selfActivate"
	GRAPH_EXPIRATION_AFTER_DURATION = "afterDuration"
)

// Default directory scope for Entra role activations (tenant-wide)
const GRAPH_DEFAULT_DIRECTORY_SCOPE string = "/"

// Microsoft Entra (AAD) authority hosts for the device code (public client) sign-in
const (
	AAD_GLOBAL_AUTHORITY_HOST = "https://login.microsoftonline.com"
	AAD_USGOV_AUTHORITY_HOST  = "https://login.microsoftonline.us"
	AAD_CN_AUTHORITY_HOST     = "https://login.partner.microsoftonline.cn"
)

// Tenant used in the authority URL when no tenant ID is configured
const AAD_DEFAULT_TENANT string = "organizations"

// Base URLs for different Azure environments
var ARM_BASE_URLS = map[string]string{
	"global": ARM_GLOBAL_BASE_URL,
	"usgov":  ARM_USGOV_BASE_URL,
	"china":  ARM_CN_BASE_URL,
}

// Microsoft Graph base URLs for different Azure environments (Entra groups/roles)
var GRAPH_BASE_URLS = map[string]string{
	"global": GRAPH_GLOBAL_BASE_URL,
	"usgov":  GRAPH_USGOV_BASE_URL,
	"china":  GRAPH_CN_BASE_URL,
}

// Microsoft Graph token scopes (.default) for different Azure environments
var GRAPH_SCOPES = map[string]string{
	"global": GRAPH_GLOBAL_BASE_URL + "/.default",
	"usgov":  GRAPH_USGOV_BASE_URL + "/.default",
	"china":  GRAPH_CN_BASE_URL + "/.default",
}

// Microsoft Entra authority hosts for different Azure environments
var AAD_AUTHORITY_HOSTS = map[string]string{
	"global": AAD_GLOBAL_AUTHORITY_HOST,
	"usgov":  AAD_USGOV_AUTHORITY_HOST,
	"china":  AAD_CN_AUTHORITY_HOST,
}
