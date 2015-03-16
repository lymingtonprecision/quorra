package ext

import (
	"errors"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/version"

	"golang.org/x/net/context"
)

const Key = "co.uk.lymingtonprecision.quorra"

type Extension struct {
	Client *govmomi.Client
	Object *types.Extension
}

func makeExtension() types.Extension {
	return types.Extension{
		Key: Key,
		Description: &types.Description{
			Label:   "Quorra",
			Summary: "A utility for managing CoreOS VMs on vSphere",
		},
		Company:                "Lymington Precision Engineers Co. Ltd.",
		LastHeartbeatTime:      time.Now(),
		Version:                version.Str(),
		ShownInSolutionManager: true,
		ExtendedProductInfo: &types.ExtExtendedProductInfo{
			CompanyUrl: "http://www.lymingtonprecision.co.uk/",
			ProductUrl: "https://github.com/lymingtonprecision/quorra",
		},
		ManagedEntityInfo: []types.ExtManagedEntityInfo{
			types.ExtManagedEntityInfo{
				Type:        "VirtualMachine",
				Description: "Virtual Machine",
			},
		},
	}
}

func Register(cl *govmomi.Client) (*Extension, error) {
	_, ok := Get(cl)
	if ok {
		return nil, errors.New("extension already registered")
	}

	re := &types.RegisterExtension{
		This:      *cl.ServiceContent.ExtensionManager,
		Extension: makeExtension(),
	}

	_, err := methods.RegisterExtension(context.TODO(), cl, re)
	if err != nil {
		return nil, err
	}

	ex, ok := Get(cl)
	if !ok {
		return nil, errors.New("failed to get extension after registering")
	}

	return ex, nil
}

func Unregister(cl *govmomi.Client) error {
	ur := &types.UnregisterExtension{
		This:         *cl.ServiceContent.ExtensionManager,
		ExtensionKey: Key,
	}

	_, err := methods.UnregisterExtension(context.TODO(), cl, ur)
	if err != nil {
		return err
	}

	return nil
}

func Get(cl *govmomi.Client) (*Extension, bool) {
	fe := &types.FindExtension{
		This:         *cl.ServiceContent.ExtensionManager,
		ExtensionKey: Key,
	}

	r, err := methods.FindExtension(context.TODO(), cl, fe)
	if err != nil {
		panic(err)
	}

	if r.Returnval != nil {
		return &Extension{Client: cl, Object: r.Returnval}, true
	}

	return nil, false
}
