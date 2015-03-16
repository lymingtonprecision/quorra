package folder

import (
	"errors"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

func GetVMFolder(cl *govmomi.Client, dc *govmomi.Datacenter, path string) (*govmomi.Folder, error) {
	if len(path) == 0 {
		folders, err := dc.Folders()
		if err != nil {
			return nil, err
		}

		return folders.VmFolder, nil
	}

	var dco mo.Datacenter
	err := cl.Properties(dc.Reference(), []string{"name"}, &dco)
	if err != nil {
		return nil, err
	}

	si := cl.SearchIndex()

	r, err := si.FindByInventoryPath(dco.Name + "/vm/" + path)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, errors.New("folder not found: " + path)
	}

	f, ok := r.(*govmomi.Folder)
	if !ok {
		return nil, errors.New("can't convert reference to folder")
	}

	return f, nil
}
