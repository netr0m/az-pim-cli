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

// Base URL for the Azure RBAC (Governance Role) API (Entra Groups and Entra Roles)
const AZ_RBAC_BASE_URL string = "https://api.azrbac.mspim.azure.com"

// Base path for the Azure Resource Manager PIM API
const ARM_BASE_PATH string = "providers/Microsoft.Authorization"

// Base path for the Azure RBAC (Governance Role) API
const AZ_RBAC_BASE_PATH = "api/v2/privilegedAccess"

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
