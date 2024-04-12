/*
Copyright © 2023 netr0m <netr0m@pm.me>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "az-pim-cli",
	Short: "A utility to list and activate Azure AD PIM roles from the CLI",
	Long: `az-pim-cli is a utility that allows the user to list and activate eligible role assignments
	from Azure Entra ID Privileged Identity Management (PIM) directly from the command line`,
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.az-pim-cli.yaml)")
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
	if err := vpr.ReadInConfig(); err == nil {
	}

	bindFlags(rootCmd, vpr)
}

func bindFlags(cmd *cobra.Command, vpr *viper.Viper) {
	cmd.Flags().VisitAll(func(flg *pflag.Flag) {
		// Replace hyphens
		configName := strings.ReplaceAll(flg.Name, "-", "")

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !flg.Changed && vpr.IsSet(configName) {
			val := vpr.Get(configName)
			cmd.Flags().Set(flg.Name, fmt.Sprintf("%v", val))
		}

	})
}
