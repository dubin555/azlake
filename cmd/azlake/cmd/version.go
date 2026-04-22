package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dubin555/azlake/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("azlake %s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
