# Azure PIM CLI
*Azure Privileged Identity Management Command Line Interface*

[![Go Reference](https://pkg.go.dev/badge/github.com/netr0m/az-pim-cli.svg)](https://pkg.go.dev/github.com/netr0m/az-pim-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/netr0m/az-pim-cli)](https://goreportcard.com/report/github.com/netr0m/az-pim-cli)

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
      --cloud string    Which Azure environment to use ('global', 'usgov', 'china') (default "global")
  -c, --config string   config file (default is $HOME/.az-pim-cli.yaml)
      --debug           Enable debug logging
  -h, --help            help for az-pim-cli

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

```bash
$ az-pim-cli list roles
```

<details>
<summary>Example</summary>

```bash
# List eligible Entra role assignments
$ az-pim-cli list roles
== my-entra-id-role ==
         - Owner
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
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Contributor scope=S100-Example-Subscription reason="" ticketNumber="" ticketSystem="" duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Contributor scope=S100-Example-Subscription status=Provisioned

# Activate a specific role ('Owner') for a resource with the prefix 's100'
$ az-pim-cli activate resource --prefix s100 --role owner
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=S100-Example-Subscription reason="" ticketNumber="" ticketSystem="" duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=S100-Example-Subscription status=Provisioned

# Activate a resource role and specify a ticket number for the activation
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --ticket-system Jira --ticket-number T-1337
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=S100-Example-Subscription reason="" ticketNumber=T-1337 ticketSystem=Jira duration=480 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=S100-Example-Subscription status=Provisioned

# Activate a resource role and specify the start time for the activation. Uses the local timezone.
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --start-time 14:30
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=S100-Example-Subscription reason="" ticketNumber=T-1337 ticketSystem=Jira duration=480 startDateTime=2024-11-20T14:30:00+01:00
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=S100-Example-Subscription status=Provisioned

# Activate a resource role and specify the start time and start date for the activation. Uses the local timezone.
$ az-pim-cli activate resource --name S100-Example-Subscription --role Owner --start-date 31/12/2024 --start-time 09:30
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=S100-Example-Subscription reason="" ticketNumber=T-1337 ticketSystem=Jira duration=480 startDateTime=2024-12-31T09:30:00+01:00
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=S100-Example-Subscription status=Provisioned
```

</details>

#### Groups
> Activate [groups](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadgroup)

```bash
$ az-pim-cli activate group
```

<details>
<summary>Example</summary>

> :information_source: See examples under [Activate - Azure resources](#azure-resources-1) for additional parameters.

```bash
# Activate the first matching role for the group 'my-entra-id-group'
$ az-pim-cli activate group --name my-entra-id-group --duration 5
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=my-entra-id-group reason="" ticketNumber="" ticketSystem="" duration=5 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned subStatus=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=my-entra-id-group status=Active
```

</details>

#### Entra roles
> Activate [entra roles](https://portal.azure.com/#view/Microsoft_Azure_PIMCommon/ActivationMenuBlade/~/aadmigratedroles)

```bash
$ az-pim-cli activate role
```
<details>
<summary>Example</summary>

> :information_source: See examples under [Activate - Azure resources](#azure-resources-1) for additional parameters.

```bash
# Activate the first matching role for the Entra role 'my-entra-id-role'
$ az-pim-cli activate role --name my-entra-id-role --duration 5
time=2024-11-20T08:08:08.534+01:00 level=INFO msg="Requesting activation" role=Owner scope=my-entra-id-role reason="" ticketNumber="" ticketSystem="" duration=5 startDateTime=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="The role assignment request was successful" status=Provisioned subStatus=""
time=2024-11-20T08:08:20.129+01:00 level=INFO msg="Request completed" role=Owner scope=my-entra-id-role status=Active
```

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
```

#### Environment variables
You may also define these configuration options as environment variables by prefixing any global variable with `PIM_`.

```bash
export PIM_DURATION=30
export PIM_CLOUD=global
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
