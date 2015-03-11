package ext

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ext"
)

type register struct {
	FlagSet  *flag.FlagSet
	Force    *bool
	CertFile *string
	KeyFile  *string
}

func (cmd *register) CommandLine() string {
	return "[--force|-f] [CERT_FILE_PATH] [KEY_FILE_PATH]"
}

func (cmd *register) Summary() string {
	return "Registers Quorra as a vSphere vCenter plugin"
}

func (cmd *register) Description() string {
	return fmt.Sprintf(
		`%s

Prints information regarding the registered extension. If the extension
does not already exist, or '--force' is specified, then a new certificate
will be generated for authenticating extension operations, with the
resulting key pair saved to CERT_FILE_PATH and KEY_FILE_PATH or otherwise
printed to stdout.

**Both parts of the certificate must be provided in quorra's
configuration to perform any subsequent extension related operations.**
`,
		cmd.Summary(),
	)
}

func (cmd *register) setFlags(args []string) error {
	fs := flag.NewFlagSet("register flags", flag.ContinueOnError)
	cmd.FlagSet = fs

	cmd.Force = fs.Bool("force", false, "force registration")
	fs.BoolVar(cmd.Force, "f", false, "force registration")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() > 0 {
		cmd.CertFile = &fs.Args()[0]
	} else {
		cmd.CertFile = nil
	}

	if fs.NArg() > 1 {
		cmd.KeyFile = &fs.Args()[1]
	} else {
		cmd.KeyFile = nil
	}

	return nil
}

func (cmd *register) Run(cl *govmomi.Client, c *config.Config, args []string) (err error) {
	var ex *ext.Extension
	var exists = false

	if err = cmd.setFlags(args); err != nil {
		return
	}

	ex, exists = ext.Get(cl)
	if !exists {
		if ex, err = ext.Register(cl); err != nil {
			return
		}
	}

	fmt.Printf("Registered\n\n%+v\n", *ex.Object)

	if !exists || *cmd.Force {
		fmt.Printf("\nGenerating new certificate...\n\n")
		return cmd.replaceCertificate(ex)
	}

	return nil
}

func (cmd *register) replaceCertificate(ex *ext.Extension) error {
	cert, err := ex.ReplaceCertificate()
	if err != nil {
		return err
	}

	if err := printKeyFile("public", cmd.CertFile, cert.WritePublicKey); err != nil {
		return err
	}

	return printKeyFile("private", cmd.KeyFile, cert.WritePrivateKey)
}

func printKeyFile(t string, dst *string, writer func(out io.Writer) error) (err error) {
	var out = os.Stdout

	if dst != nil && len(*dst) > 0 {
		out, err = os.Create(*dst)
		if err != nil {
			return
		}
		defer out.Close()

		fmt.Printf("Writing %s key to: %s\n", t, *dst)
	}

	return writer(out)
}

type unregister struct{}

func (cmd *unregister) Summary() string {
	return "Un-registers Quorra from the vSphere vCenter plugins"
}

func (cmd *unregister) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	err := ext.Unregister(cl)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully un-registered extension\n")

	return nil
}

func init() {
	cli.RegisterCommand([]string{"register"}, &register{})
	cli.RegisterCommand([]string{"unregister"}, &unregister{})
}
