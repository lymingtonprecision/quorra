package cert

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"

	"github.com/lymingtonprecision/quorra/config"
)

func FromConfig(c *config.Config) (*Certificate, error) {
	if len(c.ExtCert) == 0 {
		return nil, errors.New("no certificate specified in quorra configuration")
	}

	if len(c.ExtKey) == 0 {
		return nil, errors.New("no private key specified in quorra configuration")
	}

	cf, err := readFileOrReturn(c.ExtCert)
	if err != nil {
		return nil, err
	}

	pk, err := readFileOrReturn(c.ExtKey)
	if err != nil {
		return nil, err
	}

	kp, err := tls.X509KeyPair(cf, pk)
	if err != nil {
		return nil, err
	}

	cert := Certificate(kp)

	return &cert, nil
}

func readFileOrReturn(pathOrValue string) ([]byte, error) {
	f, err := ioutil.ReadFile(pathOrValue)

	switch err {
	case nil:
		return f, nil
	case os.ErrNotExist:
		return []byte(pathOrValue), nil
	default:
		return []byte{}, err
	}
}
