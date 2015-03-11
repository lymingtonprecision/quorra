package ext

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"regexp"

	"github.com/vmware/govmomi"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/cert"
	"github.com/lymingtonprecision/quorra/config"

	"golang.org/x/net/context"
)

func (ext *Extension) TlsClient(c *config.Config) (*govmomi.Client, error) {
	insecure := usingInsecureClientConnection(ext.Client)

	sc := newTunnelClient(insecure)

	ct, err := cert.FromConfig(c)
	if err != nil {
		return nil, err
	}

	setTransportClientCertAndProxy(sc, ct, vCenterProxy(ext.Client), insecure)

	cl := govmomi.Client(*ext.Client)
	cl.Client = sc

	if err := loginWithCertificate(&cl); err != nil {
		return nil, err
	}

	return &cl, nil
}

func newTunnelClient(insecure bool) *soap.Client {
	tunnel, _ := url.Parse("https://sdkTunnel:8089")
	return soap.NewClient(*tunnel, insecure)
}

func usingInsecureClientConnection(gcl *govmomi.Client) bool {
	t := gcl.Client.Transport.(*http.Transport)
	return t.TLSClientConfig.InsecureSkipVerify
}

func vCenterProxy(gcl *govmomi.Client) func(*http.Request) (*url.URL, error) {
	proxy := url.URL(gcl.Client.URL())
	hostPort := regexp.MustCompile(":\\d+$")

	proxy.Scheme = "http"
	proxy.Host = hostPort.ReplaceAllString(proxy.Host, ":80")

	return http.ProxyURL(&proxy)
}

func setTransportClientCertAndProxy(
	sc *soap.Client,
	ct *cert.Certificate,
	proxy func(*http.Request) (*url.URL, error),
	insecure bool,
) {
	t := sc.Client.Transport.(*http.Transport)
	t.Proxy = proxy
	t.TLSClientConfig.InsecureSkipVerify = insecure
	t.TLSClientConfig.Certificates = []tls.Certificate{tls.Certificate(*ct)}
}

func loginWithCertificate(gcl *govmomi.Client) error {
	_, err := methods.LoginExtensionByCertificate(
		context.TODO(),
		gcl,
		&types.LoginExtensionByCertificate{
			This:         *gcl.ServiceContent.SessionManager,
			ExtensionKey: Key,
			Locale:       "en",
		},
	)

	return err
}
