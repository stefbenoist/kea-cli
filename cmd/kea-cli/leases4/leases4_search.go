package leases4

import (
	"context"

	"github.com/spf13/cobra"

	"kea-cli/api"
	"kea-cli/cmd/kea-cli/internal/common"
)

func init() {
	var (
		hwAddress, hostname, clientID string
	)
	var searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search all IPv4 leases matching a given feature.",
		Long: `Search all IPv4 leases matching a given feature. This can be a specified hardware address, client ID or hostname.
The corresponding flags "hw_address", "hostname" and "client_id" are mutually exclusive.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				common.ErrExit(err)
			}
			return search(client, hwAddress, hostname, clientID)
		},
	}
	searchCmd.Flags().StringVarP(&hwAddress, "hw_address", "", "", "Search by hardware address.")
	searchCmd.Flags().StringVarP(&hostname, "hostname", "h", "", "Search by hostname.")
	searchCmd.Flags().StringVarP(&clientID, "client_id", "", "", "Search by client ID.")
	searchCmd.MarkFlagsMutuallyExclusive("hw_address", "client_id", "hostname")
	leases4Cmd.AddCommand(searchCmd)
}

func search(httpClient api.HTTPClient, hwAddress, hostname, clientID string) error {
	ctx := context.Background()

	var criteria, value string
	if len(hostname) > 0 {
		criteria = "hostname"
		value = hostname
	} else if len(hwAddress) > 0 {
		criteria = "hw-address"
		value = hwAddress
	} else {
		criteria = "client-id"
		value = clientID
	}
	leases, err := httpClient.Lease4().Search(ctx, criteria, value)
	if err != nil {
		return err
	}
	displayLeases(leases)
	return nil
}
