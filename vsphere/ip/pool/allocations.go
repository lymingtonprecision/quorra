package pool

import (
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/vsphere/ext"

	"golang.org/x/net/context"
)

func (pool *IpPool) AllocateIpv4Address(name string) (string, error) {
	ippm := pool.ipPoolManager()

	r, err := methods.AllocateIpv4Address(
		context.TODO(),
		pool.Client,
		&types.AllocateIpv4Address{
			This:         *ippm,
			Dc:           *pool.Datacenter,
			PoolId:       pool.Object.Id,
			AllocationId: name,
		},
	)
	if err != nil {
		return "", err
	}

	return r.Returnval, nil
}

func (pool *IpPool) ReleaseAllocation(name string) error {
	ippm := pool.ipPoolManager()

	_, err := methods.ReleaseIpAllocation(
		context.TODO(),
		pool.Client,
		&types.ReleaseIpAllocation{
			This:         *ippm,
			Dc:           *pool.Datacenter,
			PoolId:       pool.Object.Id,
			AllocationId: name,
		},
	)

	return err
}

func (pool *IpPool) AllocatedIpv4Addresses() ([]types.IpPoolManagerIpAllocation, error) {
	ippm := pool.ipPoolManager()

	aa, err := methods.QueryIPAllocations(
		context.TODO(),
		pool.Client,
		&types.QueryIPAllocations{
			This:         *ippm,
			Dc:           *pool.Datacenter,
			PoolId:       pool.Object.Id,
			ExtensionKey: ext.Key,
		},
	)
	if err != nil {
		return nil, err
	}

	return aa.Returnval, nil
}
