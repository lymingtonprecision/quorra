package pool

import (
	"errors"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"golang.org/x/net/context"
)

func GetAll(cl *govmomi.Client, dc string) ([]IpPool, error) {
	dcr, err := cl.SearchIndex().FindByInventoryPath(dc)
	if err != nil {
		return nil, err
	}
	dcmo := dcr.Reference()

	ippm := cl.ServiceContent.IpPoolManager

	r, err := methods.QueryIpPools(
		context.TODO(),
		cl,
		&types.QueryIpPools{
			This: *ippm,
			Dc:   dcmo,
		},
	)
	if err != nil {
		return nil, err
	}

	pools := []IpPool{}
	for _, pool := range r.Returnval {
		t := types.IpPool(pool)
		pools = append(pools, IpPool{Client: cl, Datacenter: &dcmo, Object: &t})
	}

	return pools, nil
}

func Get(cl *govmomi.Client, dc string, name string) (*IpPool, error) {
	pools, err := GetAll(cl, dc)
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		if strings.EqualFold(name, pool.Object.Name) {
			return &pool, nil
		}
	}

	return nil, errors.New("ip pool not found")
}
