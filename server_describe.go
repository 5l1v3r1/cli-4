package cli

import (
	"errors"
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func newServerDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SERVER",
		Short:                 "Describe a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerDescribe),
	}
	return cmd
}

func runServerDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	server, _, err := cli.Client().Server.Get(cli.Context, id)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", server.ID)
	fmt.Printf("Name:\t\t%s\n", server.Name)
	fmt.Printf("Status:\t\t%s\n", server.Status)
	fmt.Printf("Created:\t%s (%s)\n", server.Created, humanize.Time(server.Created))

	fmt.Printf("Server Type:\t%s (ID: %d)\n", server.ServerType.Name, server.ServerType.ID)
	fmt.Printf("  ID:\t\t%d\n", server.ServerType.ID)
	fmt.Printf("  Name:\t\t%s\n", server.ServerType.Name)
	fmt.Printf("  Description:\t%s\n", server.ServerType.Description)
	fmt.Printf("  Cores:\t%d\n", server.ServerType.Cores)
	fmt.Printf("  Memory:\t%v GB\n", server.ServerType.Memory)
	fmt.Printf("  Disk:\t\t%d GB\n", server.ServerType.Disk)
	fmt.Printf("  Storage Type:\t%s\n", server.ServerType.StorageType)

	fmt.Printf("Public Net:\n")
	fmt.Printf("  IPv4:\n")
	fmt.Printf("    IP:\t\t%s\n", server.PublicNet.IPv4.IP)
	fmt.Printf("    Blocked:\t%s\n", yesno(server.PublicNet.IPv4.Blocked))
	fmt.Printf("    DNS:\t%s\n", server.PublicNet.IPv4.DNSPtr)
	fmt.Printf("  IPv6:\n")
	fmt.Printf("    IP:\t\t%s\n", server.PublicNet.IPv6.IP)
	fmt.Printf("    Blocked:\t%s\n", yesno(server.PublicNet.IPv6.Blocked))
	fmt.Printf("  Floating IPs:\n")
	if len(server.PublicNet.FloatingIPs) > 0 {
		for _, floatingIP := range server.PublicNet.FloatingIPs {
			fmt.Printf("    - ID: %d\n", floatingIP.ID)
		}
	} else {
		fmt.Printf("    No Floating IPs\n")
	}

	fmt.Printf("Traffic:\n")
	fmt.Printf("  Outgoing:\t%v\n", humanize.Bytes(server.OutgoingTraffic))
	fmt.Printf("  Ingoing:\t%v\n", humanize.Bytes(server.IngoingTraffic))
	fmt.Printf("  Included:\t%v\n", humanize.Bytes(server.IncludedTraffic))

	if server.BackupWindow != "" {
		fmt.Printf("Backup Window:\t%s\n", server.BackupWindow)
	} else {
		fmt.Printf("Backup Window:\tBackups disabled\n")
	}

	if server.RescueEnabled {
		fmt.Printf("Rescue System:\tenabled\n")
	} else {
		fmt.Printf("Rescue System:\tdisabled\n")
	}

	fmt.Printf("ISO:\n")
	if server.ISO != nil {
		fmt.Printf("  ID:\t\t%d\n", server.ISO.ID)
		fmt.Printf("  Name:\t\t%s\n", server.ISO.Name)
		fmt.Printf("  Description:\t%s\n", server.ISO.Description)
		fmt.Printf("  Type:\t\t%s\n", server.ISO.Type)
	} else {
		fmt.Printf("  No ISO attached\n")
	}

	return nil
}