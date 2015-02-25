package cli

import (
	"github.com/lymingtonprecision/quorra/config"
	"github.com/vmware/govmomi"
)

type Command interface {
	Run(cl *govmomi.Client, c *config.Config) error
}

var commands = map[string]Command{}
var aliases = map[string]string{}

func Register(name string, c Command) {
	commands[name] = c
}

func Alias(name string, alias string) {
	aliases[alias] = name
}

func RealCommandName(alias string) (string, bool) {
	name, ok := aliases[alias]
	if !ok {
		name = alias
	}

	if _, exists := commands[name]; exists {
		return name, true
	}

	return "", false
}
