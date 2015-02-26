package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type HasSummary interface {
	Summary() string
}

type HasCommandLine interface {
	CommandLine() string
}

type HasDescription interface {
	Description() string
}

func PrintProgramHelp() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))

	fmt.Fprintf(os.Stderr, "\nAvailable commands:\n\n")

	printCommands(rootCommands.Commands, []string{})

	fmt.Fprintf(
		os.Stderr,
		"\nFor help with a command run %s help <command>\n",
		filepath.Base(os.Args[0]))
}

func printCommands(cm map[string]CommandMap, path []string) {
	keys := make([]string, len(cm))
	i := 0
	for k, _ := range cm {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		c, _ := cm[k]

		if c.This != nil {
			cp := strings.TrimSpace(strings.Join(append(path, k), " "))

			if s, ok := c.This.(HasSummary); ok {
				fmt.Fprintf(os.Stderr, "%-25s %s\n", cp, s.Summary())
			} else {
				fmt.Fprintf(os.Stderr, "%*s\n", cp)
			}
		} else if len(c.Commands) > 0 {
			printCommands(c.Commands, append(path, k))
		}
	}
}

func PrintInvalidCommand(args []string) {
	fmt.Fprintf(os.Stderr, "Error: no commands matched the command line %s", args)
	fmt.Fprintf(os.Stderr, "\n\n")
	PrintProgramHelp()
}

func PrintCommandHelp(cmd *PreparedCommand) {
	fmt.Fprintf(
		os.Stderr,
		"Usage of %s %s",
		filepath.Base(os.Args[0]),
		cmd.Name)

	if cl, ok := cmd.Source.(HasCommandLine); ok {
		fmt.Fprintf(os.Stderr, " %s", cl.CommandLine())
	}
	fmt.Fprintf(os.Stderr, "\n\n")

	if d, ok := cmd.Source.(HasDescription); ok {
		fmt.Fprintf(os.Stderr, "%s\n", d.Description())
	}
}

func HelpWith(args []string) {
	if len(args) == 0 {
		PrintProgramHelp()
	} else {
		c, ok := FindCommand(args)

		if ok {
			PrintCommandHelp(c)
		} else {
			PrintInvalidCommand(args)
		}
	}
}
