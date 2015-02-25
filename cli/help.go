package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	cmds := []string{}
	for name := range commands {
		cmds = append(cmds, name)
	}

	sort.Strings(cmds)

	for _, name := range cmds {
		if s, ok := commands[name].(HasSummary); ok {
			fmt.Fprintf(os.Stderr, " %-25s %s\n", name, s.Summary())
		} else {
			fmt.Fprintf(os.Stderr, " %s\n", name)
		}
	}

	fmt.Fprintf(
		os.Stderr,
		"\nFor help with a command run %s help <command>\n",
		filepath.Base(os.Args[0]))
}

func PrintInvalidCommand(name string) {
	fmt.Fprintf(os.Stderr, "Error: '%s' not recognized as a valid command", name)
	fmt.Fprintf(os.Stderr, "\n\n")
	PrintProgramHelp()
}

func PrintCommandHelp(name string) {
	cmd := commands[name]

	fmt.Fprintf(
		os.Stderr,
		"Usage of %s %s",
		filepath.Base(os.Args[0]),
		name)

	if cl, ok := cmd.(HasCommandLine); ok {
		fmt.Fprintf(os.Stderr, " %s", cl.CommandLine())
	}
	fmt.Fprintf(os.Stderr, "\n\n")

	if d, ok := cmd.(HasDescription); ok {
		fmt.Fprintf(os.Stderr, "%s\n", d.Description())
	}
}

func HelpWith(args []string) {
	if len(args) == 0 {
		PrintProgramHelp()
	} else {
		name, ok := RealCommandName(args[0])

		if ok {
			PrintCommandHelp(name)
		} else {
			PrintInvalidCommand(args[0])
		}
	}
}
