package pool

import (
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"golang.org/x/net/context"
)

func (pool *IpPool) Destroy(force bool) error {
	ippm := pool.ipPoolManager()

	_, err := methods.DestroyIpPool(
		context.TODO(),
		pool.Client,
		&types.DestroyIpPool{
			This:  *ippm,
			Dc:    *pool.Datacenter,
			Id:    pool.Object.Id,
			Force: force,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
