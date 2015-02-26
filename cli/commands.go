package cli

import (
	"strings"

	"github.com/lymingtonprecision/quorra/config"
	"github.com/vmware/govmomi"
)

type Command interface {
	Run(cl *govmomi.Client, c *config.Config, args []string) error
}

type PreparedCommand struct {
	Source Command
	Name   string
	Run    func(cl *govmomi.Client, co *config.Config) error
}

type CommandMap struct {
	This     Command
	Commands map[string]CommandMap
}

func (cm *CommandMap) Register(path []string, c Command) {
	var (
		m  = *cm
		ok bool
	)

	for i := 0; i < len(path)-1; i++ {
		if _, exists := m.Commands[path[i]]; !exists {
			m.Commands[path[i]] = CommandMap{Commands: map[string]CommandMap{}}
		}

		m, ok = m.Commands[path[i]]
		if !ok {
			panic("missing path element")
		}
	}

	m.Commands[path[len(path)-1]] = CommandMap{This: c}
}

func (cm *CommandMap) Find(args []string, path string) (*PreparedCommand, bool) {
	if len(args) == 0 {
		return nil, false
	}

	if c, exists := cm.Commands[args[0]]; exists {
		path = strings.TrimSpace(strings.Join([]string{path, args[0]}, " "))

		if sc, exists := c.Find(args[1:], path); exists {
			return sc, true
		}

		if c.This != nil {
			pc := PreparedCommand{
				Source: c.This,
				Name:   path,
				Run: func(cl *govmomi.Client, co *config.Config) error {
					return c.This.Run(cl, co, args[1:])
				},
			}

			return &pc, true
		}
	}

	return nil, false
}

var rootCommands = CommandMap{Commands: map[string]CommandMap{}}

func FindCommand(args []string) (*PreparedCommand, bool) {
	return rootCommands.Find(args, "")
}

func RegisterCommand(path []string, c Command) {
	rootCommands.Register(path, c)
}
