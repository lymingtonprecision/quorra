package version

import (
	"fmt"

	"github.com/vmware/govmomi"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
)

const (
	Major  = 0
	Minor  = 0
	Patch  = 1
	Branch = ""
)

func Str() string {
	mmp := fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

	if len(Branch) > 0 {
		return mmp + "-" + Branch
	}

	return mmp
}

type version struct{}

func (cmd *version) Summary() string {
	return "prints the version of quorra being used"
}

func (cmd *version) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	fmt.Printf("%s\n", Str())
	return nil
}

func init() {
	cli.RegisterCommand([]string{"version"}, &version{})
}
