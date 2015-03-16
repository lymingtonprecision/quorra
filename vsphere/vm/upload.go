package vm

import (
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/ova"
)

type vmUpload struct {
	url  *url.URL
	item types.OvfFileItem
	ch   chan progress.Report
}

func (u *vmUpload) Sink() chan<- progress.Report {
	return u.ch
}

func uploadVMFiles(
	cl *govmomi.Client,
	lease *govmomi.HttpNfcLease,
	info *types.HttpNfcLeaseInfo,
	spec *types.OvfCreateImportSpecResult,
	src *ova.Ova,
) error {
	items, err := itemList(cl, info, spec)
	if err != nil {
		return err
	}

	u := itemUploader{client: cl.Client, lease: lease, uploads: items, src: src}
	defer u.Done()

	err = u.Upload()
	if err != nil {
		return err
	}

	return nil
}

func itemList(
	cl *govmomi.Client,
	info *types.HttpNfcLeaseInfo,
	spec *types.OvfCreateImportSpecResult,
) ([]vmUpload, error) {
	var items []vmUpload

	for _, dev := range info.DeviceUrl {
		fi, found := findFileItem(spec, &dev)
		if !found {
			continue
		}

		url, err := cl.Client.ParseURL(dev.Url)
		if err != nil {
			return nil, err
		}

		items = append(items, vmUpload{url: url, item: *fi, ch: make(chan progress.Report)})
	}

	return items, nil
}

func findFileItem(spec *types.OvfCreateImportSpecResult, dev *types.HttpNfcLeaseDeviceUrl) (*types.OvfFileItem, bool) {
	for _, item := range spec.FileItem {
		if dev.ImportKey == item.DeviceId {
			return &item, true
		}
	}

	return nil, false
}
