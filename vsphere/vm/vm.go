package vm

import (
	"fmt"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ext"
)

type VirtualMachine struct {
	Client *govmomi.Client
	Object *govmomi.VirtualMachine

	mo mo.VirtualMachine
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

func Find(cl *govmomi.Client, c *config.Config, name string) (*VirtualMachine, error) {
	vms, err := FindAll(cl, c)
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		if strings.EqualFold(vm.Name(), name) {
			return &vm, nil
		}
	}

	return nil, fmt.Errorf("virtual machine '%s' not found", name)
}

func (vm *VirtualMachine) config() *types.VirtualMachineConfigInfo {
	if vm.mo.Config != nil {
		return vm.mo.Config
	}

	vm.Client.Properties(vm.Object.Reference(), []string{"config"}, &vm.mo)

	return vm.mo.Config
}

func (vm *VirtualMachine) Name() string {
	return vm.config().Name
}

func (vm *VirtualMachine) IsManagedByQuorra() bool {
	c := vm.config()

	if c.ManagedBy == nil {
		return false
	}

	return c.ManagedBy.ExtensionKey == ext.Key
}
