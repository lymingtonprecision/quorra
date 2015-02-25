package vm

import (
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"

	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/tag"
)

type VirtualMachine struct {
	Client *govmomi.Client
	Object *govmomi.VirtualMachine
}

func FromReference(cl *govmomi.Client, ref govmomi.Reference) (*VirtualMachine, bool) {
	if vm, ok := ref.(*govmomi.VirtualMachine); ok {
		return &VirtualMachine{Client: cl, Object: vm}, true
	}

	return nil, false
}

func FindAll(cl *govmomi.Client, c *config.Config) ([]VirtualMachine, error) {
	var vms []VirtualMachine

	r, err := cl.SearchIndex().FindByInventoryPath(c.VMStorePath())
	if err != nil {
		panic(err)
	}

	folder, ok := r.(*govmomi.Folder)
	if !ok {
		err := fmt.Errorf(
			"expected %s to be a Folder, got %s",
			c.VMStorePath(),
			r.Reference().Type,
		)

		return vms, err
	}

	items, err := folder.Children()
	if err != nil {
		return vms, err
	}

	for _, ch := range items {
		vm, ok := FromReference(cl, ch)

		if ok && vm.IsManagedByQuorra() {
			vms = append(vms, *vm)
		}
	}

	return vms, nil
}

func (vm *VirtualMachine) IsManagedByQuorra() bool {
	var o mo.VirtualMachine
	vm.Client.Properties(vm.Object.Reference(), []string{"customValue"}, &o)
	return tag.IsManagedEntity(vm.Client, o.ManagedEntity)
}
