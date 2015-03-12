package pool

import (
	"errors"
	"flag"
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ip/pool"
)

type rm struct {
	FlagSet *flag.FlagSet
	Force   *bool
}

func (cmd *rm) CommandLine() string {
	return "[NAME] [--force|-f]"
}

func (cmd *rm) Summary() string {
	return "Removes an IP Address Pool from the datacenter"
}

func (cmd *rm) setFlags(args []string) error {
	fs := flag.NewFlagSet("ip pool rm flags", flag.ContinueOnError)
	cmd.FlagSet = fs

	cmd.Force = fs.Bool("force", false, "force removal")
	fs.BoolVar(cmd.Force, "f", false, "force removal")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return nil
}

func (cmd *rm) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	if len(args) == 0 {
		return errors.New("no pool specified")
	}

	if err := cmd.setFlags(args[1:]); err != nil {
		return err
	}

	pl, err := pool.Get(cl, c.Datacenter, args[0])
	if err != nil {
		return err
	}

	if err := pl.Destroy(*cmd.Force); err != nil {
		return err
	}

	fmt.Printf("Removed IP Pool '%s'\n", pl.Object.Name)

	return nil
}

func init() {
	cli.RegisterCommand([]string{"ip", "pool", "rm"}, &rm{})
}
