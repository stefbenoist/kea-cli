package leases4

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"kea-cli/api"
	"kea-cli/cmd/kea-cli/internal/common"
)

func init() {
	var ipAddress, hwAddress, clientID string
	var subnetID int
	var delCmd = &cobra.Command{
		Use:   "del",
		Short: "Delete an IPv4 lease.",
		Long: `Delete an IPv4 lease: there are two different ways to use this command: 
- with the "ip" flag when the IPv4 is known and the details of the lease are not; one common use case of this type of query is to find out whether a given address is being used. 
- with the "hw_address" or "client_id" flags when the IPv4 is not known. In this cas, the "subnet_id" flag must be provided.
The corresponding flags "ip", "hw_address" and "client_id" are mutually exclusive.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				common.ErrExit(err)
			}
			return del(client, ipAddress, hwAddress, clientID, subnetID)
		},
	}
	delCmd.Flags().StringVarP(&ipAddress, "ip", "i", "", "Optional IPv4 address when known")
	delCmd.Flags().StringVarP(&hwAddress, "hw_address", "", "", "Optional hardware address when the IPv4 is not known.")
	delCmd.Flags().StringVarP(&clientID, "client_id", "", "", "Optional client ID when the IPv4 is not known.")
	delCmd.Flags().IntVarP(&subnetID, "subnet_id", "s", 0, "Optional subnet ID. Mandatory when hardware address or client ID is provided.")
	delCmd.MarkFlagsMutuallyExclusive("hw_address", "client_id", "ip")
	leases4Cmd.AddCommand(delCmd)
}

func del(httpClient api.HTTPClient, ipAddress, hwAddress, clientID string, subnetID int) error {
	ctx := context.Background()

	if err := httpClient.Lease4().Del(ctx, handleLeaseIdentifier(ipAddress, hwAddress, clientID, subnetID)); err != nil {
		return err
	}
	fmt.Printf("Lease for ipv4:%s is successfully deleted\n", ipAddress)
	return nil
}
