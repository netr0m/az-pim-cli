# Azure PIM CLI
*Azure Privileged Identity Management Command Line Interface*

[![Go Reference](https://pkg.go.dev/badge/github.com/netr0m/az-pim-cli.svg)](https://pkg.go.dev/github.com/netr0m/az-pim-cli)

`az-pim-cli` eases the process of listing and activating Azure PIM roles by allowing activation via the command line. Authentication is handled with the `azure.identity` library by utilizing the `AzureCLICredential` method.
It currently supports ['azure resources'](#azure-resources), ['groups'](#groups), and ['entra roles'](#entra-roles)

## Install
### Install with `go install`
```bash
$ go install github.com/netr0m/az-pim-cli@latest
```

### Clone and build yourself
```bash
# Clone the git repo
$ git clone https://github.com/netr0m/az-pim-cli.git

# Navigate into the repo directory and build
$ cd az-pim-cli
$ go build

# Move the az-pim-cli binary into your path
$ mv ./az-pim-cli /usr/local/bin
```

## Configuration
In addition to supporting environment variables and command line arguments, the script also supports certain config parameters stored in a file. By default, the script will try to look for a YAML config file at `$HOME/.az-pim-cli.yaml`, but you may also override the config file to use by supplying the `--config` flag.
See [Configuration options](#configuration-options) for more details

### Prerequisites
This tool depends on [`az-cli`](https://learn.microsoft.com/en-us/cli/azure/) for authentication. Please ensure that you've authenticated with your Azure tenant by running the command `az login`. A new browser window will open, asking you to authenticate. This should only be necessary to do once.

## Usage
```bash
$ az-pim-cli --help
az-pim-cli is a utility that allows the user to list and activate eligible role assignments
        from Azure Entra ID Privileged Identity Management (PIM) directly from the command line.

Usage:
  az-pim-cli [command]

Available Commands:
  activate    Send a request to Azure PIM to activate a role assignment
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Query Azure PIM for eligible role assignments
  version     Display the version of az-pim-cli

Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
      --debug           Enable debug logging
  -h, --help            help for az-pim-cli

Use "az-pim-cli [command] --help" for more information about a command.

```

### List eligible role assignments

#### Azure resources
> List [azure resources](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/azurerbac)

```bash
$ az-pim-cli list resources --help
Query Azure PIM for eligible resource assignments (azure resources)

Usage:
  az-pim-cli list resource [flags]

Aliases:
  resource, r, res, resource, resources, sub, subs, subscriptions

Flags:
  -h, --help   help for resource

Global Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
      --debug           Enable debug logging

```

#### Groups
> List [groups](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup)
>
> :warning: Requires an access token with the appropriate scope. See [Token for Entra ID Groups](#token-for-entra-id-groups) for more details.

```bash
$ az-pim-cli list groups --help
Query Azure PIM for eligible group assignments

Usage:
  az-pim-cli list group [flags]

Aliases:
  group, g, grp, groups

Flags:
  -h, --help           help for group
  -t, --token string   An access token for the PIM Groups API (required). Consult the README for more information.

Global Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
      --debug           Enable debug logging

```

#### Entra roles
> List [entra roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles)
>
> :warning: Requires an access token with the appropriate scope. See [Token for Entra ID Groups and Roles](#token-for-entra-id-groups-and-roles) for more details.

```bash
$ az-pim-cli list roles --help
Query Azure PIM for eligible Entra role assignments

Usage:
  az-pim-cli list role [flags]

Aliases:
  role, rl, role, roles

Flags:
  -h, --help           help for role
  -t, --token string   An access token for the PIM 'Entra Roles' and 'Groups' API (required). Consult the README for more information.

Global Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
      --debug           Enable debug logging

```

### Activate a role

#### Azure resources
> Activate [azure resources](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/azurerbac)

```bash
$ az-pim-cli activate resource --help
Sends a request to Azure PIM to activate the given resource (azure resources)

Usage:
  az-pim-cli activate resource [flags]

Aliases:
  resource, r, res, resource, resources, sub, subs, subscriptions

Flags:
  -h, --help   help for resource

Global Flags:
  -c, --config string          config file (default is $HOME/.az-pim-cli.yaml)
      --debug                  Enable debug logging
      --dry-run                Display the resource that would be activated, without requesting the activation
  -d, --duration int           Duration in minutes that the role should be activated for (default 480)
  -n, --name string            The name of the resource to activate
  -p, --prefix string          The name prefix of the resource to activate (e.g. 'S399'). Alternative to 'name'.
      --reason string          Reason for the activation (default "config")
  -r, --role string            Specify the role to activate, if multiple roles are found for a resource (e.g. 'Owner' and 'Contributor')
  -T, --ticket-number string   Ticket number for the activation
      --ticket-system string   Ticket system for the activation

```

#### Groups
> Activate [groups](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup)
>
> :warning: Requires an access token with the appropriate scope. See [Token for Entra ID Groups](#token-for-entra-id-groups) for more details.

```bash
$ az-pim-cli activate group --help
Sends a request to Azure PIM to activate the given group

Usage:
  az-pim-cli activate group [flags]

Aliases:
  group, g, grp, groups

Flags:
  -h, --help           help for group
  -t, --token string   An access token for the PIM 'Entra Roles' and 'Groups' API (required). Consult the README for more information.

Global Flags:
  -c, --config string          config file (default is $HOME/.az-pim-cli.yaml)
      --debug                  Enable debug logging
      --dry-run                Display the resource that would be activated, without requesting the activation
  -d, --duration int           Duration in minutes that the role should be activated for (default 480)
  -n, --name string            The name of the resource to activate
  -p, --prefix string          The name prefix of the resource to activate (e.g. 'S399'). Alternative to 'name'.
      --reason string          Reason for the activation (default "config")
  -r, --role string            Specify the role to activate, if multiple roles are found for a resource (e.g. 'Owner' and 'Contributor')
  -T, --ticket-number string   Ticket number for the activation
      --ticket-system string   Ticket system for the activation

```

#### Entra roles
> Activate [entra roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles)
>
> :warning: Requires an access token with the appropriate scope. See [Token for Entra ID Groups and Roles](#token-for-entra-id-groups-and-roles) for more details.

```bash
$ az-pim-cli activate role --help
go run main.go activate role --help
Sends a request to Azure PIM to activate the given Entra role

Usage:
  az-pim-cli activate role [flags]

Aliases:
  role, rl, role, roles

Flags:
  -h, --help           help for role
  -t, --token string   An access token for the PIM 'Entra Roles' and 'Groups' API (required). Consult the README for more information.

Global Flags:
  -c, --config string          config file (default is $HOME/.az-pim-cli.yaml)
      --debug                  Enable debug logging
      --dry-run                Display the resource that would be activated, without requesting the activation
  -d, --duration int           Duration in minutes that the role should be activated for (default 480)
  -n, --name string            The name of the resource to activate
  -p, --prefix string          The name prefix of the resource to activate (e.g. 'S399'). Alternative to 'name'.
      --reason string          Reason for the activation (default "config")
  -r, --role string            Specify the role to activate, if multiple roles are found for a resource (e.g. 'Owner' and 'Contributor')
  -T, --ticket-number string   Ticket number for the activation
      --ticket-system string   Ticket system for the activation

```

### Examples
#### Azure resources
```bash
# List eligible Azure resource role assignments
$ az-pim-cli list resources
== S100-Example-Subscription ==
         - Contributor
         - Owner
== S1337-Another-Subscription ==
         - Contributor

# Activate the first matching role for a resource with the prefix 'S100'
$ az-pim-cli activate resource --prefix S100
2024/05/31 15:05:25 Activating role 'Contributor' for resource 'S100-Example-Subscription' with reason 'config' (ticket:  [])
2024/05/31 15:05:34 The role 'Contributor' in 'S100-Example-Subscription' is now Provisioned

# Activate a specific role ('Owner') for a resource with the prefix 's100'
$ az-pim-cli activate resource --prefix s100 --role owner
2024/05/31 15:06:25 Activating role 'Owner' for resource 'S100-Example-Subscription' with reason 'config' (ticket:  [])
2024/05/31 15:06:34 The role 'Owner' in 'S100-Example-Subscription' is now Provisioned

# Activate a resource role and specify a ticket number for the activation
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --ticket-system Jira --ticket-number T-1337
2024/05/31 15:06:25 Activating role 'Owner' for resource 'S100-Example-Subscription' with reason 'config' (ticket: T-1337 [Jira])
2024/05/31 15:06:34 The role 'Owner' in 'S100-Example-Subscription' is now Provisioned
```

#### Groups
```bash
# List eligible group assignments
$ az-pim-cli list groups
== my-entra-id-group ==
         - Owner

# Activate the first matching role for the group 'my-entra-id-group'
$ az-pim-cli activate group --name my-entra-id-group --duration 5
2024/05/31 15:00:10 Activating role 'Owner' for group 'my-entra-id-group' with reason 'config' (ticket:  [])
2024/05/31 15:00:23 The role 'Owner' for group 'my-entra-id-group' is now Active
```

#### Entra roles
```bash
# List eligible Entra role assignments
$ az-pim-cli list roles
== my-entra-id-role ==
         - Owner
         - Contributor

# Activate the first matching role for the Entra role 'my-entra-id-role'
$ az-pim-cli activate role --name "my-entra-id-role" --duration 30
2024/05/31 15:00:10 Activating role 'Owner' for Entra role 'my-entra-id-role' with reason 'config' (ticket:  [])
2024/05/31 15:00:23 The role 'Owner' for Entra role 'my-entra-id-role' is now Active

# Activate nominated role for the Entra role 'my-entra-id-role'
$ az-pim-cli activate role --name "my-entra-id-role" --role "Contributor" --duration 30
2024/05/31 15:00:10 Activating role 'Contributor' for Entra role 'my-entra-id-role' with reason 'config' (ticket:  [])
2024/05/31 15:00:23 The role 'Contributor' for Entra role 'my-entra-id-role' is now Active
```

### Configuration options

- `token`: The Bearer token to use for authorization when requesting the Azure PIM Groups endpoint, i.e. listing/activating Azure PIM Groups and Entra Roles

#### YAML file
You may define configuration options in a YAML file.
By default, the program will use the file ~/.az-pim-cli.yaml ($HOME/.az-pim-cli.yaml), if present. You may override this path with the command line flag `--config [PATH]`.

```bash
$ cat ~/.az-pim-cli.yaml
token: eyJ0[...]
reason: static-reason
ticketSystem: System
ticketNumber: T-1337
duration: 5
```

#### Environment variables
You may also define these configuration options as environment variables by prefixing any global variable with `PIM_`.

```bash
export PIM_TOKEN=eyJ0[...]

```

### Token for Entra ID Groups and Roles
Due to limitations with authorization for Azure PIM, this software may only acquire a token authorized for listing and activating ['Azure resources' roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/azurerbac).
In order to list or activate ['Entra groups'](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup) and ['Entra roles'](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles), you must acquire a token from an authenticated browser session. This token will have a limited lifetime, which means you'll likely have to perform this step each time you wish to activate or list Entra groups.

To acquire the token, do the following:
1. Navigate to ['Microsoft Entra Privileged Identity Management > Activate > Groups'](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup) or ['Microsoft Entra Privileged Identity Management > Activate > Microsoft Entra roles'](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles)
2. Open *DevTools* (`CTRL+Shift+I`), and locate a request to `https://api.azrbac.mspim.azure.com/api/v2/privilegedAccess/aadGroups/roleAssignments` or `https://api.azrbac.mspim.azure.com/api/v2/privilegedAccess/aadroles/roleAssignments`
    - If no such request can be seen, press the "Refresh" button above the table to issue a new request
    - In *DevTools*, the "File" attribute should start with "roleAssignments"
3. In *DevTools*, under the "Headers" tab for the given request, copy the value of the `Authorization` header, which should start with "Bearer eyJ0[...]"
4. Remove the prefix "Bearer" from the value, resulting in "eyJ0[...]"
5. Set an environment variable or config file value according to the description in [Configuration options](#configuration-options), e.g.
  ```
  PIM_TOKEN=eyJ0[...]
  ```
6. You may now, and for the duration of the token's lifetime, list and activate 'Entra groups' and 'Entra roles' using this tool

### Troubleshooting

To ease the process of troubleshooting, you can add the flag `--debug` to enable debug logging.

> :warning: Debug logs contain sensitive information. Take care to sensor any sensitive data before sharing the output.

```bash
$ az-pim-cli activate role --name my-entra-id-role --duration 5 --debug
```

## Contributing

Want to contribute to the project? There are a few things you need to know.

See [CONTRIBUTING](./CONTRIBUTING.md) to get started
