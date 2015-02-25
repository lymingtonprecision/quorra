package tag

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const ManagedTagName string = "quorra.ManagedByQuorra"

func IsManagedEntity(cl *govmomi.Client, e mo.ManagedEntity) bool {
	td, err := tagDef(cl, ManagedTagName)
	if err != nil {
		panic(err)
	}

	for _, cv := range e.CustomValue {
		cfv, ok := cv.(*types.CustomFieldStringValue)

		if ok && cfv.Key == td.Key && cfv.Value == "true" {
			return true
		}
	}

	return false
}

func tagDef(cl *govmomi.Client, tagName string) (*types.CustomFieldDef, error) {
	cfm := cl.ServiceContent.CustomFieldsManager
	var cfs mo.CustomFieldsManager

	err := cl.Properties(*cfm, []string{"field"}, &cfs)
	if err != nil {
		return nil, err
	}

	for _, cf := range cfs.Field {
		if cf.Name == tagName {
			return &cf, nil
		}
	}

	r, err := methods.AddCustomFieldDef(
		cl,
		&types.AddCustomFieldDef{This: *cfm, Name: tagName},
	)
	if err != nil {
		return nil, err
	}

	return &r.Returnval, nil
}
