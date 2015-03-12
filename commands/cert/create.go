package cert

import (
	"fmt"
	"io"
	"os"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cert"
	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
)

type create struct{}

func (cmd *create) CommandLine() string {
	return "[CERT_FILE_PATH] [KEY_FILE_PATH]"
}

func (cmd *create) Summary() string {
	return "Creates a new X509 key pair, " +
		"saving them to the specified locations or printing to stdout"
}

func (cmd *create) Description() string {
	return fmt.Sprintf(
		`%s

Takes two (optional) arguments specifying the file paths to save the certificate
and private key file to:

    cert create cert.pem cert.key

Will write the public certificate to 'cert.pem' and the private key file
to 'cert.key' in the current directory.

If only one file path is provided the private key will be printed to stdout.
If no paths are provided both the public and private keys will be printed.
`,
		cmd.Summary(),
	)
}

func (cmd *create) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	cert, err := cert.Create()
	if err != nil {
		return err
	}

	writeCertFiles(cert, args)

	return nil
}

func writeCertFiles(crt *cert.Certificate, args []string) error {
	var certOut, keyOut io.Writer
	var err error

	if len(args) >= 1 {
		certOut, err = os.Create(args[0])
		if err != nil {
			return err
		}
		defer certOut.(io.Closer).Close()

		fmt.Printf("Writing public cert to: %s\n\n", args[0])
	} else {
		certOut = os.Stdout
	}

	if len(args) >= 2 {
		keyOut, err = os.Create(args[1])
		if err != nil {
			return err
		}
		defer keyOut.(io.Closer).Close()

		fmt.Printf("Writing private key to: %s\n\n", args[1])
	} else {
		keyOut = os.Stdout
	}

	crt.WritePublicKey(certOut)
	crt.WritePrivateKey(keyOut)

	return nil
}

func init() {
	cli.RegisterCommand([]string{"cert", "create"}, &create{})
}
