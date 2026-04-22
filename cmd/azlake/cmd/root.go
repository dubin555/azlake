package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dubin555/azlake/pkg/logging"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "azlake",
	Short: "azlake is a data lake management platform",
	Long:  `azlake - a data lake management platform based on lakeFS concepts`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is azlake.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("azlake")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("/etc/azlake")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found; ignore
		_ = err
	}
	logging.SetLevel("info")
}
