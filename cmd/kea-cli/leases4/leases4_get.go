package leases4

import (
	"context"

	"github.com/spf13/cobra"

	"kea-cli/api"
	"kea-cli/cmd/kea-cli/internal/common"
)

func init() {
	var ipAddress, hwAddress, clientID string
	var subnetID int
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get an IPv4 lease.",
		Long: `Get an IPv4 lease: there are two different ways to use this command: 
- with the "ip" flag when the IPv4 is known and the details of the lease are not; one common use case of this type of query is to find out whether a given address is being used. 
- with the "hw_address" or "client_id" flags when the IPv4 is not known. In this cas, the "subnet_id" flag must be provided.
The corresponding flags "ip", "hw_address" and "client_id" are mutually exclusive.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				common.ErrExit(err)
			}
			return get(client, ipAddress, hwAddress, clientID, subnetID)
		},
	}
	getCmd.Flags().StringVarP(&ipAddress, "ip", "i", "", "Optional IPv4 address when known")
	getCmd.Flags().StringVarP(&hwAddress, "hw_address", "", "", "Optional hardware address when the IPv4 is not known.")
	getCmd.Flags().StringVarP(&clientID, "client_id", "", "", "Optional client ID when the IPv4 is not known.")
	getCmd.Flags().IntVarP(&subnetID, "subnet_id", "s", 0, "Optional subnet ID. Mandatory when hardware address or client ID is provided.")
	getCmd.MarkFlagsMutuallyExclusive("hw_address", "client_id", "ip")
	leases4Cmd.AddCommand(getCmd)
}

func get(httpClient api.HTTPClient, ipAddress, hwAddress, clientID string, subnetID int) error {
	ctx := context.Background()

	lease, err := httpClient.Lease4().Get(ctx, handleLeaseIdentifier(ipAddress, hwAddress, clientID, subnetID))
	if err != nil {
		return err
	}
	displayLeases([]api.Lease{*lease})
	return nil
}

func handleLeaseIdentifier(ipAddress, hwAddress, clientID string, subnetID int) *api.LeaseIdentifier {
	identifier := api.LeaseIdentifier{}
	if len(ipAddress) > 0 {
		identifier.IpAddress = ipAddress
	} else if len(hwAddress) > 0 {
		identifier.IdentifierType = "hw-address"
		identifier.IdentifierValue = hwAddress
		identifier.SubnetID = subnetID
	} else if len(clientID) > 0 {
		identifier.IdentifierType = "client-id"
		identifier.IdentifierValue = clientID
		identifier.SubnetID = subnetID
	}
	return &identifier
}
