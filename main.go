package main

import (
	"os"

	"github.com/lymingtonprecision/quorra/cli"
	_ "github.com/lymingtonprecision/quorra/commands"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
