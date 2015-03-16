package vm

import (
	"fmt"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/soap"

	"github.com/lymingtonprecision/quorra/cli/progresslogger"
	"github.com/lymingtonprecision/quorra/ova"
)

type itemUploader struct {
	uploads []vmUpload
	src     *ova.Ova

	lease  *govmomi.HttpNfcLease
	client *soap.Client

	pos   int64
	total int64

	done chan struct{}
	wg   sync.WaitGroup
}

func (u *itemUploader) Upload() error {
	u.done = make(chan struct{})

	// create progress update watchers for all the uploads
	for _, upload := range u.uploads {
		u.total += upload.item.Size
		go u.waitForUploadProgress(upload)
	}

	// kick start the lease renewal
	u.wg.Add(1)
	go u.keepLeaseRefreshed()

	// perform the uploads one by one
	for _, upload := range u.uploads {
		if err := u.performUpload(upload); err != nil {
			return err
		}
	}

	return nil
}

func (u *itemUploader) performUpload(upload vmUpload) error {
	f, err := u.src.GetFile(upload.item.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	logger := progresslogger.NewProgressLogger(fmt.Sprintf("* Uploading %s...", path.Base(upload.item.Path)))
	defer logger.Wait()

	opts := soap.Upload{
		ContentLength: upload.item.Size,
		Progress:      progress.Tee(&upload, logger),
		Method:        "POST",
		Type:          "application/x-vnd.vmware-streamVmdk",
	}

	if upload.item.Create {
		opts.Method = "PUT"
		opts.Type = ""
		opts.Headers = map[string]string{
			"Overwrite": "t",
		}
	}

	return u.client.Upload(f, upload.url, &opts)
}

func (u *itemUploader) waitForUploadProgress(upload vmUpload) {
	var pos, total int64

	total = upload.item.Size

	for {
		select {
		case <-u.done:
			return
		case p, ok := <-upload.ch:
			if ok && p.Error() != nil {
				return
			}

			// we've hit the last element
			if !ok {
				atomic.AddInt64(&u.pos, total-pos)
				return
			}

			// approximate progress in bytes
			x := int64(float32(total) * (p.Percentage() / 100.0))
			atomic.AddInt64(&u.pos, x-pos)
			pos = x
		}
	}
}

func (u *itemUploader) keepLeaseRefreshed() {
	defer u.wg.Done()

	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-u.done:
			return
		case <-tick.C:
			pcnt := int(float32(100*atomic.LoadInt64(&u.pos)) / float32(u.total))
			err := u.lease.HttpNfcLeaseProgress(pcnt)
			if err != nil {
				fmt.Printf("from lease updater: %s\n", err)
			}
		}
	}
}

func (u *itemUploader) Done() {
	close(u.done)
	u.wg.Wait()
}
