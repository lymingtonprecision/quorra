package cli

import (
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere"
)

func Run(args []string) int {
	if len(args) == 0 {
		PrintProgramHelp()
		return 1
	}

	if args[0] == "help" {
		HelpWith(args[1:])
		return 1
	}

	name, ok := RealCommandName(args[0])

	if !ok {
		PrintInvalidCommand(args[0])
		return 1
	}

	c, err := config.FromDefaultSources()
	if err != nil {
		panic(err)
	}

	cl, err := vsphere.NewClient(c)
	if err != nil {
		panic(err)
	}

	if err := commands[name].Run(cl, c); err != nil {
		panic(err)
	}

	return 0
}
