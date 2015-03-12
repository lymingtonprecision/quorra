package pool

import (
	"errors"
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ip/pool"
)

type info struct{}

func (cmd *info) CommandLine() string {
	return "[NAME]"
}

func (cmd *info) Summary() string {
	return "Prints information about IP pool NAME"
}

func (cmd *info) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	if len(args) == 0 {
		return errors.New("no pool name supplied")
	}

	pl, err := pool.Get(cl, c.Datacenter, args[0])
	if err != nil {
		return err
	}

	printPoolsSummary([]pool.IpPool{*pl})

	aa, err := pl.AllocatedIpv4Addresses()
	if err != nil {
		return err
	}

	if len(aa) > 0 {
		fmt.Println("\nAllocated Addresses:")

		for _, addr := range aa {
			fmt.Printf("%s\t%s\n", addr.AllocationId, addr.IpAddress)
		}
	}

	return nil
}

func init() {
	cli.RegisterCommand([]string{"ip", "pool", "info"}, &info{})
}
