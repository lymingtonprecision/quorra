package vm

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"
)

func (vm *VirtualMachine) AddNetwork(n *govmomi.Network) error {
	return vm.Object.AddDevice(
		&types.VirtualVmxnet3{
			types.VirtualVmxnet{
				types.VirtualEthernetCard{
					VirtualDevice: types.VirtualDevice{
						UnitNumber: -1,
						Backing: &types.VirtualEthernetCardNetworkBackingInfo{
							VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{
								DeviceName: n.Name(),
							},
							InPassthroughMode: false,
						},
						Connectable: &types.VirtualDeviceConnectInfo{
							StartConnected:    true,
							AllowGuestControl: true,
							Connected:         true,
						},
					},
					AddressType:      string(types.VirtualEthernetCardMacTypeGenerated),
					WakeOnLanEnabled: true,
				},
			},
		},
	)
}
