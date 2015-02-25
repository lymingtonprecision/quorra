package vsphere

import (
	"net/url"
	"regexp"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/config"
)

func NewClient(config *config.Config) (*govmomi.Client, error) {
	u, err := vCenterURL(config)
	if err != nil {
		return nil, err
	}

	c, err := govmomi.NewClient(*u, true)
	if err != nil {
		return nil, err
	}

	ui := url.UserPassword(config.Username, config.Password)
	if err := c.Login(*ui); err != nil {
		return nil, err
	}

	return c, err
}

func vCenterURL(config *config.Config) (*url.URL, error) {
	addr := config.ApiUrl
	if b, _ := regexp.MatchString("^\\w+://", addr); !b {
		addr = "https://" + addr
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.Path == "" {
		u.Path = "/sdk"
	}

	if u.User != nil {
		u.User = nil
	}

	return u, nil
}
