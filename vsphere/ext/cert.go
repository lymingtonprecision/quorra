package ext

import (
	"bytes"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/cert"

	"golang.org/x/net/context"
)

func (ext *Extension) UploadCertificate(c *cert.Certificate) error {
	var b bytes.Buffer

	if err := c.WritePublicKey(&b); err != nil {
		return err
	}

	_, err := methods.SetExtensionCertificate(
		context.TODO(),
		ext.Client,
		&types.SetExtensionCertificate{
			This:           *ext.Client.ServiceContent.ExtensionManager,
			ExtensionKey:   Key,
			CertificatePem: b.String(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (ext *Extension) ReplaceCertificate() (*cert.Certificate, error) {
	c, err := cert.Create()
	if err != nil {
		return nil, err
	}

	err = ext.UploadCertificate(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
