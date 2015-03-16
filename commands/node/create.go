package node

import (
	"flag"
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/ova"
	"github.com/lymingtonprecision/quorra/vsphere/vm"
)

type create struct {
	Flags *flag.FlagSet
	Host  *string
	Ova   *string
}

func (cmd *create) CommandLine() string {
	return "[--host|-h HOST] [--ova OVA_PATH]"
}

func (cmd *create) Summary() string {
	return "Creates a new CoreOS node"
}

func (cmd *create) Description() string {
	return fmt.Sprintf(
		`%s

Creates either a standalone node or joins the existing nodes in a cluster.
Starts the VM on either the specified HOST or the default host from the
configuration file.

Uses the OVA file specified in the configuration or from OVA_PATH if
provided.
		`,
		cmd.Summary(),
	)
}

func (cmd *create) setFlags(c *config.Config, args []string) error {
	fs := flag.NewFlagSet("node create flags", flag.ContinueOnError)

	cmd.Flags = fs

	cmd.Host = fs.String("host", "", "host on which to run the VM")
	fs.StringVar(cmd.Host, "h", "", "host on which to run the VM")

	cmd.Ova = fs.String("ova", "", "path to OVA file to import")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if len(*cmd.Ova) == 0 {
		cmd.Ova = &c.Default.OVA
	}

	return nil
}

func (cmd *create) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	if err := cmd.setFlags(c, args); err != nil {
		return err
	}

	ova, err := ova.Open(*cmd.Ova)
	if err != nil {
		return err
	}
	defer ova.Close()

	ref, err := cmd.getObjectReferences(cl, c)
	if err != nil {
		return err
	}

	_, err = vm.CreateFromOva(cl, c, ref, ova)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *create) getObjectReferences(cl *govmomi.Client, c *config.Config) (*config.References, error) {
	return c.GetVMReferences(cl, config.Overrides{Host: *cmd.Host})
}

func init() {
	cli.RegisterCommand([]string{"node", "create"}, &create{})
}
