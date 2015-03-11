package cert

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

func (cert *Certificate) WritePublicKey(out io.Writer) error {
	return pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]})
}

func (cert *Certificate) privKeyPemBlock() (pb *pem.Block, err error) {
	switch k := cert.PrivateKey.(type) {
	case *rsa.PrivateKey:
		pb = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		pb, err = &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
		if err != nil {
			return nil, err
		}
	default:
		err = errors.New(
			fmt.Sprintf("unable to marshal %T private key", k),
		)
	}

	return
}

func (cert *Certificate) WritePrivateKey(out io.Writer) error {
	b, err := cert.privKeyPemBlock()
	if err != nil {
		return err
	}

	return pem.Encode(out, b)
}
