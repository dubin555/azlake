package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd handles azlakectl configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create or update the azlakectl configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			fmt.Printf("Config file: %s\n", cfgFile)
		} else {
			fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())
		}
		fmt.Printf("Server endpoint: %s\n", cfg.Server.EndpointURL)
		fmt.Printf("Access Key ID: %s\n", cfg.Credentials.AccessKeyID)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

