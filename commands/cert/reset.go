package cert

import (
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cert"
	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
)

type reset struct{}

func (cmd *reset) CommandLine() string {
	return "[CERT_FILE_PATH] [KEY_FILE_PATH]"
}

func (cmd *reset) Summary() string {
	return "Equivilant to calling `create` then `upload`"
}

func (cmd *reset) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	crt, err := cert.Create()
	if err != nil {
		return err
	}

	err = writeCertFiles(crt, args)
	if err != nil {
		return err
	}

	err = uploadCert(cl, crt)
	if err != nil {
		return err
	}

	fmt.Println("\nCertificate reset")

	return nil
}

func init() {
	cli.RegisterCommand([]string{"cert", "reset"}, &reset{})
}
