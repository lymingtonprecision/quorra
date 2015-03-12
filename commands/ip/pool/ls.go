package pool

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"text/tabwriter"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ip/pool"
)

type ls struct{}

func (cmd *ls) Summary() string {
	return "Provides a summary of all the IP pools configured in the datacenter"
}

func (cmd *ls) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	pools, err := pool.GetAll(cl, c.Datacenter)
	if err != nil {
		return err
	}

	printPoolsSummary(pools)

	return nil
}

func printPoolsSummary(pools []pool.IpPool) {
	rangeSize := regexp.MustCompile("\\d+$")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 2, ' ', 0)

	fmt.Fprintln(
		w,
		"ID\tName\tRange\tAllocated\tAvailable\tNetmask\tGateway\t",
	)

	for _, pl := range pools {
		s, _ := strconv.Atoi(rangeSize.FindString(pl.Object.Ipv4Config.Range))
		aa, _ := pl.AllocatedIpv4Addresses()

		fmt.Fprintf(
			w,
			"%d\t%s\t%s\t%d\t%d\t%s\t%s\t\n",
			pl.Object.Id,
			pl.Object.Name,
			pl.Object.Ipv4Config.Range,
			len(aa),
			s-len(aa),
			pl.Object.Ipv4Config.Netmask,
			pl.Object.Ipv4Config.Gateway,
		)
	}
	w.Flush()
}

func init() {
	cli.RegisterCommand([]string{"ip", "pool", "ls"}, &ls{})
}
