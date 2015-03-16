package vm

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/namesgenerator"
	"github.com/lymingtonprecision/quorra/ova"
)

func CreateFromOva(cl *govmomi.Client, c *config.Config, ref *config.References, src *ova.Ova) (*VirtualMachine, error) {
	descriptor, err := src.ReadDescriptor()
	if err != nil {
		return nil, err
	}

	attempts := int64(0)
	maxAttempts := namesgenerator.Possibilities()

begin:
	name := namesgenerator.RandomName(0)

	spec, err := CreateImportSpec(cl, c, ref, descriptor, name)
	if err != nil {
		return nil, err
	}

	lease, err := ref.ResourcePool.ImportVApp(spec.ImportSpec, ref.Folder, ref.Host)
	if err != nil {
		return nil, err
	}

	info, err := lease.Wait()
	if isVMAlreadyExistsError(err) {
		if attempts < maxAttempts {
			attempts++
			goto begin
		} else {
			return nil, errors.New("could not generate a unique VM name")
		}
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("Creating '%s'...\n\n", name)

	err = uploadVMFiles(cl, lease, info, spec, src)
	if err != nil {
		return nil, err
	}

	err = lease.HttpNfcLeaseComplete()
	if err != nil {
		return nil, err
	}

	m, err := Find(cl, c, name)
	if err != nil {
		return nil, err
	}

	fmt.Printf("* Adding public network interface... ")
	if err := m.AddNetwork(ref.PublicNetwork); err != nil {
		return nil, err
	}
	fmt.Printf("OK\n")

	fmt.Printf("\nVM successfully created!\n")

	return m, nil
}

func isVMAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	re := regexp.MustCompile("(?i)the name '[A-Za-z0-9\\-]+' already exists\\.")
	return re.MatchString(err.Error())
}
