package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"kea-cli/api"
	"kea-cli/cmd/kea-cli/leases4"
)

var cfgFile string

var clientConfig = new(api.Configuration)

// Both version and commit are filled by Makefile
var version string
var commit string

func main() {
	initFlags()
	rootCmd.AddCommand(leases4.Init(clientConfig))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kea-cli",
	Short: "kea-cli is a Kea REST API client",
	Long:  `kea-cli is a Kea REST API client`,

	DisableAutoGenTag: true,
	SilenceUsage:      true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := getClientConfig(clientConfig)
		if err != nil {
			fmt.Print(err)
			return err
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Print(err)
		}
	},
}
