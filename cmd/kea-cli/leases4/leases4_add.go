package leases4

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"kea-cli/api"
	"kea-cli/cmd/kea-cli/internal/common"
)

func init() {
	var subnetID int
	var addCmd = &cobra.Command{
		Use:   "add <IPv4> <MAC address>",
		Short: "Add a new IPv4 lease.",
		Long:  `Add a new IPv4 lease by providing the IPv4 and the hardware (MAC) address. The "subnet_id" flag is optional. If not specified, Kea tries to determine the value by running a subnet-selection procedure.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				common.ErrExit(err)
			}
			return add(client, args[0], args[1], subnetID)
		},
	}
	addCmd.Flags().IntVarP(&subnetID, "subnet_id", "s", 0, "Optional subnet ID. If not specified, Kea tries to determine the value by running a subnet-selection procedure.")
	leases4Cmd.AddCommand(addCmd)
}

func add(httpClient api.HTTPClient, ipAddress, hwAddress string, subnetID int) error {
	ctx := context.Background()
	if err := httpClient.Lease4().Add(ctx, ipAddress, hwAddress, subnetID); err != nil {
		return err
	}
	fmt.Printf("Lease for IPv4:%s is successfully added\n", ipAddress)
	return nil
}
