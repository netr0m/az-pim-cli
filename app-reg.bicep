extension graphV1
targetScope = 'resourceGroup'

@description('Application registration display name.')
param appRegistrationName string = 'az-pim-cli'

resource app 'Microsoft.Graph/applications@v1.0' = {
  uniqueName: appRegistrationName
  displayName: appRegistrationName
  signInAudience: 'AzureADMyOrg'
  // Allow public client flows
  isFallbackPublicClient: true
  requiredResourceAccess: [
		{
      // Microsoft Graph
			resourceAppId: '00000003-0000-0000-c000-000000000000'
			resourceAccess: [
				{
          // PrivilegedAssignmentSchedule.ReadWrite.AzureADGroup - Delegated
					id: '06dbc45d-6708-4ef0-a797-f797ee68bf4b'
					type: 'Scope'
				}
				{
          // PrivilegedEligibilitySchedule.Read.AzureADGroup - Delegated
					id: '8f44f93d-ecef-46ae-a9bf-338508d44d6b'
					type: 'Scope'
				}
				{
          // RoleAssignmentSchedule.ReadWrite.Directory - Delegated
          id: '8c026be3-8e26-4774-9372-8d5d6f21daff'
					type: 'Scope'
				}
				{
          // RoleEligibilitySchedule.Read.Directory - Delegated
					id: 'eb0788c2-6d4e-4658-8c9e-c0fb8053f03d'
					type: 'Scope'
				}
				{
          // User.Read - Delegated
					id: 'e1fe6dd8-ba31-4d61-89e7-88639da4683d'
					type: 'Scope'
				}
			]
		}
	]
  web: {
    homePageUrl: 'https://github.com/netr0m/az-pim-cli'
  }
}

resource sp 'Microsoft.Graph/servicePrincipals@v1.0' = {
  appId: app.appId
  displayName: appRegistrationName
}

output appId string = app.appId
