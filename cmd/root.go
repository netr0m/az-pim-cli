/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/netr0m/az-pim-cli/pkg/common"
	"github.com/netr0m/az-pim-cli/pkg/pim"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	debugLogging        bool
	cfgFile             string
	azureEnv            string
	clientID            string
	tenantID            string
	useDeviceCode       bool
	AzureClientInstance pim.AzureClient
)

var rootCmd = &cobra.Command{
	Use:   "az-pim-cli",
	Short: "A utility to list and activate Azure AD PIM roles from the CLI",
	Long: `az-pim-cli is a utility that allows the user to list and activate eligible role assignments
	from Azure Entra ID Privileged Identity Management (PIM) directly from the command line.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		armBaseURL, ok := pim.ARM_BASE_URLS[azureEnv]
		if !ok {
			fmt.Printf("Invalid value for --cloud: %q (allowed: global, usgov, china)\n", azureEnv)
			os.Exit(1)
		}
		graphBaseURL, ok := pim.GRAPH_BASE_URLS[azureEnv]
		if !ok {
			fmt.Printf("Could not find matching Microsoft Graph base URL for the environment %q\n", azureEnv)
			os.Exit(1)
		}
		graphScope, ok := pim.GRAPH_SCOPES[azureEnv]
		if !ok {
			fmt.Printf("Could not find matching Microsoft Graph scope for the environment %q\n", azureEnv)
			os.Exit(1)
		}
		authorityHost, ok := pim.AAD_AUTHORITY_HOSTS[azureEnv]
		if !ok {
			fmt.Printf("Could not find matching authority host for the environment %q\n", azureEnv)
			os.Exit(1)
		}
		authorityTenant := tenantID
		if authorityTenant == "" {
			authorityTenant = pim.AAD_DEFAULT_TENANT
		}
		AzureClientInstance = pim.AzureClient{
			ARMBaseURL:    armBaseURL,
			GraphBaseURL:  graphBaseURL,
			GraphScope:    graphScope,
			Authority:     fmt.Sprintf("%s/%s", authorityHost, authorityTenant),
			ClientID:      clientID,
			TenantID:      tenantID,
			UseDeviceCode: useDeviceCode,
		}
	},
}

// requireGraphClient ensures a custom app registration is configured before
// attempting PIM operations for groups or Entra roles (which require Microsoft
// Graph access that the Azure CLI's built-in client cannot provide).
func requireGraphClient() {
	if AzureClientInstance.ClientID == "" {
		fmt.Println("PIM for Groups and Entra roles requires a custom app registration.")
		fmt.Println("Set --client-id and --tenant-id (or PIM_CLIENTID/PIM_TENANTID, or clientid/tenantid in ~/.az-pim-cli.yaml).")
		fmt.Println("See the README for setup instructions.")
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.az-pim-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&azureEnv, "cloud", "global", "Which Azure environment to use ('global', 'usgov', 'china')")
	rootCmd.PersistentFlags().StringVar(&clientID, "client-id", "", "Client ID of the app registration used for PIM group/role activation via Microsoft Graph (or set PIM_CLIENTID / 'clientid' in config)")
	rootCmd.PersistentFlags().StringVar(&tenantID, "tenant-id", "", "Microsoft Entra tenant ID for the app registration (or set PIM_TENANTID / 'tenantid' in config)")
	rootCmd.PersistentFlags().BoolVar(&useDeviceCode, "device-code", false, "Use the device code sign-in flow instead of the interactive browser flow (for headless environments; may be blocked by Conditional Access)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	vpr := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		vpr.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".az-pim-cli" (without extension).
		vpr.AddConfigPath(home)
		vpr.SetConfigType("yaml")
		vpr.SetConfigName(".az-pim-cli")
	}

	vpr.SetEnvPrefix("PIM")
	vpr.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	vpr.AutomaticEnv()

	// If a config file is found, read it in.
	if err := vpr.ReadInConfig(); err != nil {
		// If user provided a config file
		if cfgFile != "" {
			fmt.Printf("error loading config file: %v", err)
			os.Exit(1)
		}
	}

	bindFlags(rootCmd, vpr)
	bindFlags(activateCmd, vpr)
	bindFlags(listGroupCmd, vpr)
	bindFlags(listEntraRoleCmd, vpr)
	bindFlags(activateResourceCmd, vpr)
	bindFlags(activateGroupCmd, vpr)
	bindFlags(activateEntraRoleCmd, vpr)

	common.InitLogger(debugLogging)
}

func bindFlags(cmd *cobra.Command, vpr *viper.Viper) {
	cmd.Flags().VisitAll(func(flg *pflag.Flag) {
		// Replace hyphens
		configName := strings.ReplaceAll(flg.Name, "-", "")

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !flg.Changed && vpr.IsSet(configName) {
			val := vpr.Get(configName)
			cmd.Flags().Set(flg.Name, fmt.Sprintf("%v", val)) //nolint:errcheck
		}
	})
}
