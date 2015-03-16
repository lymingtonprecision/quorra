package node

import (
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/vm"
)

type ls struct{}

func (cmd *ls) Summary() string {
	return "Lists all current nodes"
}

func (cmd *ls) Description() string {
	return cmd.Summary()
}

func (cmd *ls) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	vms, err := vm.FindAll(cl, c)
	if err != nil {
		return err
	}

	if len(vms) > 0 {
		for _, vm := range vms {
			fmt.Printf("%v\n", vm.Object)
		}
	} else {
		fmt.Printf("* * *\nNo Managed VMs found\n* * *\n")
	}

	return nil
}

func init() {
	cli.RegisterCommand([]string{"node", "ls"}, &ls{})
}
