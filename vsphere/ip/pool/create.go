package pool

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"golang.org/x/net/context"
)

type IpPool struct {
	Client     *govmomi.Client
	Datacenter *types.ManagedObjectReference
	Object     *types.IpPool
}

func Create(cl *govmomi.Client, dc string, ipp *types.IpPool) (*IpPool, error) {
	dcr, err := cl.SearchIndex().FindByInventoryPath(dc)
	if err != nil {
		return nil, err
	}
	dcmo := dcr.Reference()

	ippm := cl.ServiceContent.IpPoolManager

	_, err = methods.CreateIpPool(
		context.TODO(),
		cl,
		&types.CreateIpPool{
			This: *ippm,
			Dc:   dcmo,
			Pool: *ipp,
		},
	)
	if err != nil {
		return nil, err
	}

	return Get(cl, dc, ipp.Name)
}

func (pool *IpPool) ipPoolManager() *types.ManagedObjectReference {
	return pool.Client.ServiceContent.IpPoolManager
}
