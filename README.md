# Azure PIM CLI
*Azure Privileged Identity Management Command Line Interface*

`az-pim-cli` eases the process of listing and activating Azure PIM roles by allowing activation via the command line. Authentication is handled with the `azure.identity` library by utilizing the `AzureCLICredential` method.

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
        from Azure Entra ID Privileged Identity Management (PIM) directly from the command line

Usage:
  az-pim-cli [command]

Available Commands:
  activate    Sends a request to Azure PIM to activate the given role
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Query Azure PIM for eligible role assignments

Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
  -h, --help            help for az-pim-cli

Use "az-pim-cli [command] --help" for more information about a command.

```

### List eligible role assignments (Azure resources)
```bash
$ az-pim-cli list --help
Query Azure PIM for eligible role assignments

Usage:
  az-pim-cli list [flags]
  az-pim-cli list [command]

Aliases:
  list, l, ls

Available Commands:
  group       Query Azure PIM for eligible group assignments

Flags:
  -h, --help   help for list

Global Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)

Use "az-pim-cli list [command] --help" for more information about a command.

```

### List eligible group assignments (Entra Groups)
> :warn: Requires an access token with the appropriate scope. See [Token for Entra ID Groups](#token-for-entra-id-groups) for more details.
```bash
$ az-pim-cli list group --help
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

```

### Activate a role (Azure resources)
```bash
$ az-pim-cli activate --help
Sends a request to Azure PIM to activate the given role

Usage:
  az-pim-cli activate [flags]
  az-pim-cli activate [command]

Aliases:
  activate, a, ac, act

Available Commands:
  group       Sends a request to Azure PIM to activate the given group

Flags:
  -d, --duration int    Duration in minutes that the role should be activated for (default 480)
  -h, --help            help for activate
  -n, --name string     The name of the resource to activate
  -p, --prefix string   The name prefix of the resource to activate (e.g. 'S399'). Alternative to 'name'.
      --reason string   Reason for the activation (default "config")
  -r, --role string     Specify the role to activate, if multiple roles are found for a subscription (e.g. 'Owner' and 'Contributor')

Global Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)

Use "az-pim-cli activate [command] --help" for more information about a command.

```

### Activate a role (Entra Groups)
> :warn: Requires an access token with the appropriate scope. See [Token for Entra ID Groups](#token-for-entra-id-groups) for more details.
```bash
$ az-pim-cli activate group --help
Sends a request to Azure PIM to activate the given group

Usage:
  az-pim-cli activate group [flags]

Aliases:
  group, g, grp, groups

Flags:
  -h, --help           help for group
  -t, --token string   An access token for the PIM Groups API (required). Consult the README for more information.

Global Flags:
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
  -d, --duration int    Duration in minutes that the role should be activated for (default 480)
  -n, --name string     The name of the resource to activate
  -p, --prefix string   The name prefix of the resource to activate (e.g. 'S399'). Alternative to 'name'.
      --reason string   Reason for the activation (default "config")
  -r, --role string     Specify the role to activate, if multiple roles are found for a subscription (e.g. 'Owner' and 'Contributor')

```

### Examples
#### Azure resources
```bash
# List eligible Azure resource role assignments
$ az-pim-cli list
== S100-Example-Subscription ==
         - Contributor
         - Owner
== S1337-Another-Subscription ==
         - Contributor

# Activate the first matching role in a subscription with the prefix 'S100'
$ az-pim-cli activate --prefix S100
2024/05/31 15:05:25 Activating role 'Contributor' in subscription 'S100-Example-Subscription' with reason 'config'
2024/05/31 15:05:34 The role 'Contributor' in 'S100-Example-Subscription' is now Provisioned

# Activate a specific role ('Owner') in a subscription with the prefix 's100'
$ az-pim-cli activate --prefix s100 --role owner
2024/05/31 15:06:25 Activating role 'Owner' in subscription 'S100-Example-Subscription' with reason 'config'
2024/05/31 15:06:34 The role 'Owner' in 'S100-Example-Subscription' is now Provisioned
```

#### Entra groups
```bash
# List eligible group assignments
$ az-pim-cli list groups
== my-entra-id-group ==
         - Owner

# Activate the first matching role for the group 'my-entra-id-group'
$ az-pim-cli activate group --name my-entra-id-group --duration 5
2024/05/31 15:00:10 Activating role 'Owner' for group 'my-entra-id-group' with reason 'config'
2024/05/31 15:00:23 The role 'Owner' for group 'my-entra-id-group' is now Active
```

### Configuration options

- `token`: The Bearer token to use for authorization when requesting the Azure PIM Groups endpoint, i.e. listing/activating Azure PIM Groups

#### YAML file
You may define global configuration options in a YAML file.
By default, the program will use the file ~/.az-pim-cli.yaml ($HOME/.az-pim-cli.yaml), if present. You may override this path with the command line flag `--config [PATH]`.

```bash
$ cat ~/.az-pim-cli.yaml
token: eyJ0[...]

```

#### Environment variables
You may also define these configuration options as environment variables by prefixing any global variable with `PIM_`.

```bash
export PIM_TOKEN=eyJ0[...]

```

### Token for Entra ID Groups
Due to limitations with authorization for Azure PIM, this software may only acquire a token authorized for listing and activating ['Azure resources' roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/azurerbac).
In order to list or activate ['Entra groups'](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup), you must acquire a token from an authenticated browser session. This token will have a limited lifetime, which means you'll likely have to perform this step each time you wish to activate or list Entra groups.

To acquire the token, do the following:
1. Navigate to ['Microsoft Entra Privileged Identity Management > Activate > Groups'](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup)
2. Open *DevTools* (`CTRL+Shift+I`), and locate a request to `https://api.azrbac.mspim.azure.com/api/v2/privilegedAccess/aadGroups/roleAssignments`
    - If no such request can be seen, press the "Refresh" button above the table to issue a new request
    - In *DevTools*, the "File" attribute should start with "roleAssignments"
3. In *DevTools*, under the "Headers" tab for the given request, copy the value of the `Authorization` header, which should start with "Bearer eyJ0[...]"
4. Remove the prefix "Bearer" from the value, resulting in "eyJ0[...]"
5. Set an environment variable or config file value according to the description in [Configuration options](#configuration-options), e.g.
  ```
  PIM_TOKEN=eyJ0[...]
  ```
6. You may now, and for the duration of the token's lifetime, list and activate 'Entra groups' using this tool
