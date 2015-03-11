package ext

import (
	"bytes"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/cert"

	"golang.org/x/net/context"
)

func (ext *Extension) ReplaceCertificate() (*cert.Certificate, error) {
	var b bytes.Buffer

	cert, err := cert.Create()
	if err != nil {
		return nil, err
	}

	if err := cert.WritePublicKey(&b); err != nil {
		return nil, err
	}

	_, err = methods.SetExtensionCertificate(
		context.TODO(),
		ext.Client,
		&types.SetExtensionCertificate{
			This:           *ext.Client.ServiceContent.ExtensionManager,
			ExtensionKey:   Key,
			CertificatePem: b.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
