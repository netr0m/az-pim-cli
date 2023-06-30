/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package pim

// Base URL for the Azure PIM API
const AZ_PIM_BASE_URL string = "https://api.azrbac.mspim.azure.com"

// Base path for the Azure PIM API
const AZ_PIM_BASE_PATH string = "api/v2/privilegedAccess"

// Authority used for Azure authentication
const AZ_AUTHORITY string = "https://login.microsoftonline.com/"

// Scope used for Azure authentication (Azure PIM AppId)
const AZ_PIM_SCOPE string = "01fc33a7-78ba-4d2f-a4b7-768e336e890e/.default"
