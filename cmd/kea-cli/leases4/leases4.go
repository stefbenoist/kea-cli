package leases4

import (
	"fmt"

	"github.com/spf13/cobra"

	"kea-cli/api"
)

// ClientConfig is the client configuration resolved by cobra/viper
var clientConfig *api.Configuration

// Init creates the leases4 command
func Init(cc *api.Configuration) *cobra.Command {
	clientConfig = cc
	return leases4Cmd
}

// templatesCmd is the templates-based command
var leases4Cmd = &cobra.Command{
	Use:           "leases4",
	Aliases:       []string{"lease4", "lease"},
	Short:         "Perform commands on leases for ipv4",
	Long:          `Perform commands on leases for ipv4`,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Print(err)
		}
	},
}
