# Azure PIM CLI
*Azure Privileged Identity Management Command Line Interface*

[![Go Reference](https://pkg.go.dev/badge/github.com/netr0m/az-pim-cli.svg)](https://pkg.go.dev/github.com/netr0m/az-pim-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/netr0m/az-pim-cli)](https://goreportcard.com/report/github.com/netr0m/az-pim-cli)

`az-pim-cli` eases the process of listing and activating Azure PIM roles by allowing activation via the command line.
It currently supports ['azure resources'](#azure-resources), ['groups'](#groups), and ['entra roles'](#entra-roles)

Authentication differs by resource type:
- **Azure resources** authenticate via your existing `az login` session (`azure.identity`'s `AzureCLICredential`).
- **Entra roles and groups** authenticate with a Microsoft Entra app registration (device code sign-in) - see [Prerequisites](#prerequisites). You'll need to provide your own Entra app registration to use with this tool. See [Azure/azure-cli#22775](https://github.com/Azure/azure-cli/issues/22775) for background.

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

#### Azure resources
This tool depends on [`az-cli`](https://learn.microsoft.com/en-us/cli/azure/) for authentication. Please ensure that you've authenticated with your Azure tenant by running the command `az login`. A new browser window will open, asking you to authenticate. This should only be necessary to do once.

#### Entra roles and groups
`az-pim-cli` is an open-source tool and does not provide an app registration. To list/activate Entra roles or groups, you'll need to register an app of your own.

The first `list`/`activate` call for roles or groups opens a device code sign-in in your browser; a token is then cached locally, so this is only needed again once it expires.

##### Switching accounts
`az-pim-cli login` clears the cached token and runs a fresh device code sign-in, so you can switch which Entra account the Entra role and Entra group commands use.

`az-pim-cli logout` clears the cached token, effectively logging you out.

##### US Gov / China Cloud
`--cloud usgov`/`--cloud china` (see [Configuration options](#configuration-options)) route Entra role/group sign-in through the matching national-cloud authority (`login.microsoftonline.us`/`login.partner.microsoftonline.cn`) and Graph endpoint (`graph.microsoft.us`/`microsoftgraph.chinacloudapi.cn`).

> :warning: The Graph API calls (list/activate) are unverified - I don't have a US Gov or China tenant to test against. Please open an issue if you hit problems.

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
  login       Sign in for the Entra roles/groups commands, replacing any cached account
  logout      Clear the cached sign-in used for the Entra roles/groups commands
  version     Display the version of az-pim-cli

Flags:
      --client-id string   Client ID of your own app registration, required for Entra role/group commands (see README)
      --cloud string       Which Azure environment to use ('global', 'usgov', 'china') (default "global")
  -c, --config string      config file (default is $HOME/.az-pim-cli.yaml)
      --debug              Enable debug logging
  -h, --help               help for az-pim-cli
      --tenant-id string   Tenant ID of your own app registration, required for Entra role/group commands (see README)

Use "az-pim-cli [command] --help" for more information about a command.

```

### List eligible role assignments

#### Azure resources
> List [azure resources](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/azurerbac)

```bash
$ az-pim-cli list resources
```

<details>
<summary>Example</summary>

```bash
# List eligible Azure resource role assignments
$ az-pim-cli list resources
== S100-Example-Subscription ==
        - Contributor
        - Owner
== S1337-Another-Subscription ==
        - Contributor
```

</details>

#### Groups
> List [groups](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup)
>
> Requires `--client-id`/`--tenant-id` (or config/env) - see [Prerequisites](#entra-roles-and-groups).

```bash
$ az-pim-cli list groups
```

<details>
<summary>Example</summary>

```bash
# List eligible group assignments
$ az-pim-cli list groups
== my-entra-id-group ==
         - Owner
```

</details>

#### Entra roles
> List [entra roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles)
>
> Requires `--client-id`/`--tenant-id` (or config/env) - see [Prerequisites](#entra-roles-and-groups).

```bash
$ az-pim-cli list roles
```

<details>
<summary>Example</summary>

```bash
# List eligible Entra role assignments
$ az-pim-cli list roles
== Global Reader ==
         - Global Reader
```

</details>

### Activate a role

#### Azure resources
> Activate [azure resources](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/azurerbac)

```bash
$ az-pim-cli activate resource
```

<details>
<summary>Examples</summary>

```bash
# Activate the first matching role for a resource with the prefix 'S100'
$ az-pim-cli activate resource --prefix S100
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Contributor scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 reason="" ticketNumber="" ticketSystem="" duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Contributor scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 status=Provisioned

# Activate a specific role ('Owner') for a resource with the prefix 's100'
$ az-pim-cli activate resource --prefix s100 --role owner
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 reason="" ticketNumber="" ticketSystem="" duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 status=Provisioned

# Activate a specific role ('Owner') for a resource with a narrower scope
$ az-pim-cli activate resource --name S100-Example-Subscription --scope /subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83/resourceGroups/example-rg --role Owner
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83/resourceGroups/example-rg reason="" ticketNumber="" ticketSystem="" duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83/resourceGroups/example-rg status=Provisioned

# Activate a resource role and specify a ticket number for the activation
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --ticket-system Jira --ticket-number T-1337
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 reason="" ticketNumber=T-1337 ticketSystem=Jira duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 status=Provisioned

# Activate a resource role and specify the start time for the activation. Uses the local timezone.
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --start-time 14:30
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 reason="" ticketNumber=T-1337 ticketSystem=Jira duration=480 startDateTime=2024-11-20T14:30:00+01:00
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 status=Provisioned

# Activate a resource role and specify the start time and start date for the activation. Uses the local timezone.
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --start-date 31/12/2024 --start-time 09:30
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 reason="" ticketNumber=T-1337 ticketSystem=Jira duration=480 startDateTime=2024-12-31T09:30:00+01:00
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=/subscriptions/6d80fbe7-7508-46e8-9fbf-a9d2e0dfec83 status=Provisioned
```

</details>

#### Groups
> Activate [groups](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup)
>
> Requires `--client-id`/`--tenant-id` (or config/env) - see [Prerequisites](#entra-roles-and-groups).

```bash
$ az-pim-cli activate group
```

<details>
<summary>Example</summary>

> :information_source: See examples under [Activate - Azure resources](#azure-resources-1) for additional parameters.

```bash
# Activate the first matching role for the group 'my-entra-id-group'
$ az-pim-cli activate group --name my-entra-id-group --duration 5
time=2026-08-26T08:08:08.534+02:00 level=INFO msg="Requesting activation" role=Owner scope=my-entra-id-group reason=config ticketNumber="" ticketSystem="" duration=5 startDateTime=2026-08-26T08:08:08+02:00 cloud=global
time=2026-08-26T08:08:08.900+02:00 level=INFO msg="The role assignment request was successful" status=PendingProvisioning
time=2026-08-26T08:08:08.900+02:00 level=INFO msg="Request completed" role=Owner scope=my-entra-id-group status=PendingProvisioning
```

</details>

#### Entra roles
> Activate [entra roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles)
>
> Requires `--client-id`/`--tenant-id` (or config/env) - see [Prerequisites](#entra-roles-and-groups).

```bash
$ az-pim-cli activate role
```
<details>
<summary>Example</summary>

> :information_source: See examples under [Activate - Azure resources](#azure-resources-1) for additional parameters.

```bash
# Activate the first matching role for the Entra role 'Global Reader'
$ az-pim-cli activate role --name "Global Reader" --duration 5
time=2026-08-26T08:08:08.534+02:00 level=INFO msg="Requesting activation" role="Global Reader" scope="Global Reader" reason=config ticketNumber="" ticketSystem="" duration=5 startDateTime=2026-08-26T08:08:08+02:00 cloud=global
time=2026-08-26T08:08:08.900+02:00 level=INFO msg="The role assignment request was successful" status=Granted
time=2026-08-26T08:08:08.900+02:00 level=INFO msg="Request completed" role="Global Reader" scope="Global Reader" status=Granted
```

> :information_source: `role` and `scope` are the same value for Entra roles.

</details>

### Configuration options

#### YAML file
You may define configuration options in a YAML file.
By default, the program will use the file ~/.az-pim-cli.yaml ($HOME/.az-pim-cli.yaml), if present. You may override this path with the command line flag `--config [PATH]`.

```bash
$ cat ~/.az-pim-cli.yaml
reason: static-reason
ticketSystem: System
ticketNumber: T-1337
duration: 5
cloud: global
clientId: 00000000-0000-0000-0000-000000000001
tenantId: 00000000-0000-0000-0000-000000000002
```

#### Environment variables
You may also define these configuration options as environment variables by prefixing any global variable with `PIM_`.

```bash
export PIM_DURATION=30
export PIM_CLOUD=global
export PIM_CLIENTID=00000000-0000-0000-0000-000000000001
export PIM_TENANTID=00000000-0000-0000-0000-000000000002
```

### Troubleshooting

To ease the process of troubleshooting, you can add the flag `--debug` to enable debug logging.

> :warning: Debug logs contain sensitive information. Take care to sensor any sensitive data before sharing the output.

```bash
$ az-pim-cli activate role --name my-entra-id-role --duration 5 --debug
```

## Testing

To run the unit tests, run the following command from the project root:

```bash
$ go test -v ./...
```

## Contributing

Want to contribute to the project? There are a few things you need to know.

See [CONTRIBUTING](./CONTRIBUTING.md) to get started
