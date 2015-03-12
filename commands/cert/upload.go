package cert

import (
	"errors"
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cert"
	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ext"
)

type upload struct{}

func (cmd *upload) Summary() string {
	return "Uploads the current X509 key pair to vSphere for use" +
		" when authenticating extension actions"
}

func (cmd *upload) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	crt, err := cert.FromConfig(c)
	if err != nil {
		return err
	}

	err = uploadCert(cl, crt)
	if err != nil {
		return err
	}

	fmt.Println("Extension Certificate Uploaded")

	return nil
}

func uploadCert(cl *govmomi.Client, crt *cert.Certificate) error {
	ext, ok := ext.Get(cl)
	if !ok {
		return errors.New("extension not found, have you registered it?")
	}

	err := ext.UploadCertificate(crt)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	cli.RegisterCommand([]string{"cert", "upload"}, &upload{})
}
