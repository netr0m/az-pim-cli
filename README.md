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
  -c, --config string      config file (default is $HOME/.az-pim-cli.yaml)
  -h, --help               help for az-pim-cli

Use "az-pim-cli [command] --help" for more information about a command.

```

### List eligible role assignments
```bash
$ az-pim-cli list --help
Query Azure PIM for eligible role assignments

Usage:
  az-pim-cli list [flags]

Aliases:
  list, l, ls

Flags:
  -h, --help   help for list

Global Flags:
  -c, --config string      config file (default is $HOME/.az-pim-cli.yaml)
```

### Activate a role
```bash
$ az-pim-cli activate --help
Sends a request to Azure PIM to activate the given role

Usage:
  az-pim-cli activate [flags]

Aliases:
  activate, a, ac, act

Flags:
  -d, --duration int                 Duration in minutes that the role should be activated for (default 480)
  -h, --help                         help for activate
      --reason string                Reason for the activation (default "config")
  -r, --role-name string             Specify the role to activate, if multiple roles are found for a subscription (e.g. 'Owner' and 'Contributor')
  -s, --subscription-name string     The name of the subscription to activate
  -p, --subscription-prefix string   The name prefix of the subscription to activate (e.g. 'S399'). Alternative to 'subscription-name'.

Global Flags:
  -c, --config string      config file (default is $HOME/.az-pim-cli.yaml)
```

### Examples
```bash
# List eligible role assignments
$ az-pim-cli list
Opening in existing browser session.
== S398-XXX ==
         - Owner
         - Contributor
== S250-XXX ==
         - Contributor

# Activate the first matching role in a subscription with the prefix 's398'
$ az-pim-cli activate --subscription-prefix s398 --duration 60
Opening in existing browser session.
2023/06/30 14:27:04 Activating role 'Owner' in subscription 'S398-XXX'
2023/06/30 14:27:11 The role 'Owner' in 'S398-XXX' is now Active

# Activate a specific role ('Owner') in a subscription with the prefix 's398'
$ az-pim-cli activate -p s398 --role-name owner
```
