package vm

import (
	"errors"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ext"
)

func CreateImportSpec(cl *govmomi.Client, c *config.Config, r *config.References, descriptor string, name string) (*types.OvfCreateImportSpecResult, error) {
	s, err := cl.OvfManager().CreateImportSpec(
		descriptor,
		r.ResourcePool,
		r.Datastore,
		types.OvfCreateImportSpecParams{
			EntityName:       name,
			DiskProvisioning: "thin",
			NetworkMapping: []types.OvfNetworkMapping{
				types.OvfNetworkMapping{
					Name:    "VM Network",
					Network: r.PrivateNetwork.Reference(),
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if s.Error != nil {
		return nil, errors.New(s.Error[0].LocalizedMessage)
	}

	is := s.ImportSpec.(*types.VirtualMachineImportSpec)
	cs := is.ConfigSpec
	cs = ensureDevicesHaveUnitNumbers(cs)
	cs = addManagedBy(cs)
	cs = setMemory(cs, c.VmMemoryInMB())

	is.ConfigSpec = cs
	s.ImportSpec = is

	return s, nil
}

func ensureDevicesHaveUnitNumbers(cs types.VirtualMachineConfigSpec) types.VirtualMachineConfigSpec {
	// skip device 0
	for _, d := range cs.DeviceChange[1:] {
		n := &d.GetVirtualDeviceConfigSpec().Device.GetVirtualDevice().UnitNumber
		if *n == 0 {
			*n = -1
		}
	}

	return cs
}

func addManagedBy(cs types.VirtualMachineConfigSpec) types.VirtualMachineConfigSpec {
	cs.ManagedBy = &types.ManagedByInfo{
		ExtensionKey: ext.Key,
		Type:         "VirtualMachine",
	}

	return cs
}

func setMemory(cs types.VirtualMachineConfigSpec, mb int64) types.VirtualMachineConfigSpec {
	cs.MemoryHotAddEnabled = true
	cs.MemoryMB = mb

	return cs
}
