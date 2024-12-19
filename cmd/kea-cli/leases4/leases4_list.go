package leases4

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"kea-cli/api"
	"kea-cli/cmd/kea-cli/internal/common"
)

func init() {
	var from string
	var limit int
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List IPv4 leases.",
		Long:  `List IPv4 leases using paging mechanism.`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				common.ErrExit(err)
			}
			return list(client, from, limit)
		},
	}
	listCmd.Flags().StringVarP(&from, "from", "f", "start", "Optional IPv4 address defining from where starting the page result")
	listCmd.Flags().IntVarP(&limit, "limit", "l", api.DefaultPageLimit, "Limit number of leases result.")
	leases4Cmd.AddCommand(listCmd)
}

func list(httpClient api.HTTPClient, from string, limit int) error {
	ctx := context.Background()
	leases, err := httpClient.Lease4().List(ctx, from, limit)
	if err != nil {
		return err
	}
	fmt.Printf("From: %s\n", from)
	fmt.Printf("Limit: %d\n", limit)
	fmt.Printf("Total: %d\n", len(leases))
	displayLeases(leases)
	return nil
}

func displayLeases(leases []api.Lease) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"Client ID", "CLTT", "FQDN FWD", "FQDN REV", "Hostname", "MAC Address", "IPv4", "State", "Subnet ID", "Valid LFT"}
	t.AppendHeader(header)
	t.AppendSeparator()
	// Limit max width for columns
	width := common.GetUsableColWidth(len(header))
	columnConfigs := make([]table.ColumnConfig, 10)
	for i := 1; i <= 10; i++ {
		columnConfigs = append(columnConfigs, table.ColumnConfig{
			Number:           i,
			WidthMax:         common.GetRatio(width, 0.1),
			AlignHeader:      text.AlignCenter,
			WidthMaxEnforcer: text.WrapSoft,
		})
	}
	t.SetColumnConfigs(columnConfigs)

	for _, lease := range leases {
		t.AppendRow(table.Row{lease.ClientID, lease.Cltt, lease.FQDNFwd, lease.FQDNRev,
			lease.Hostname, lease.HwAddress, lease.IpAddress, lease.State, lease.SubnetID, lease.ValidLft})
		t.AppendSeparator()
	}
	t.Render()
}
