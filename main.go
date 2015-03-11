package main

import (
	"os"

	"github.com/lymingtonprecision/quorra/cli"
	_ "github.com/lymingtonprecision/quorra/commands/cert"
	_ "github.com/lymingtonprecision/quorra/commands/ext"
	_ "github.com/lymingtonprecision/quorra/commands/node"
	_ "github.com/lymingtonprecision/quorra/version"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
