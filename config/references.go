package config

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"

	"github.com/lymingtonprecision/quorra/vsphere/folder"
)

type References struct {
	Datacenter     *govmomi.Datacenter
	Datastore      *govmomi.Datastore
	Folder         *govmomi.Folder
	Host           *govmomi.HostSystem
	ResourcePool   *govmomi.ResourcePool
	PublicNetwork  *govmomi.Network
	PrivateNetwork *govmomi.Network
}

type Overrides struct {
	Host string
}

func fsv(str ...string) string {
	for _, s := range str {
		if len(s) > 0 {
			return s
		}
	}

	return ""
}

func (c *Config) GetVMReferences(cl *govmomi.Client, o Overrides) (*References, error) {
	var r = References{}
	var err error

	f := find.NewFinder(cl, false)

	r.Datacenter, err = f.Datacenter(c.Datacenter)
	if err != nil {
		return nil, err
	}

	f = f.SetDatacenter(r.Datacenter)

	r.Datastore, err = f.Datastore(fsv(c.VM.Datastore, c.Default.Datastore))
	if err != nil {
		return nil, err
	}
	r.Folder, err = folder.GetVMFolder(cl, r.Datacenter, fsv(c.VM.Folder, c.Default.Folder))
	if err != nil {
		return nil, err
	}
	r.Host, err = f.HostSystem("*/" + fsv(o.Host, c.Default.Host))
	if err != nil {
		return nil, err
	}
	r.ResourcePool, err = r.Host.ResourcePool()
	if err != nil {
		return nil, err
	}
	nr, err := f.Network(c.Network.Public.Name)
	if err != nil {
		return nil, err
	}
	r.PublicNetwork = nr.(*govmomi.Network)

	nr, err = f.Network(c.Network.Private.Name)
	if err != nil {
		return nil, err
	}
	r.PrivateNetwork = nr.(*govmomi.Network)

	return &r, nil
}
