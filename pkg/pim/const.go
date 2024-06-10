/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

// Base URL for the Azure Entra PIM API
const AZ_PIM_BASE_URL string = "https://management.azure.com"

// Base path for the Azure Entra PIM API
const AZ_PIM_BASE_PATH string = "providers/Microsoft.Authorization"

// Authority used for Azure authentication
const AZ_AUTHORITY string = "https://login.microsoftonline.com/"

// Scope used for Azure authentication
const AZ_PIM_SCOPE string = AZ_PIM_BASE_URL

// Default reason for role activation
const DEFAULT_REASON string = "config"

// Default duration for role activation
const DEFAULT_DURATION_MINUTES int = 480

// API version for the "role eligibility schedule instances" (i.e. eligible azure resource role assignments)
const AZ_PIM_API_VERSION string = "2020-10-01"
