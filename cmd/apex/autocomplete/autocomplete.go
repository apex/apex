// Package autocomplete generates bash auto-completion words.
package autocomplete

import (
	"fmt"
	"io/ioutil"

	"github.com/tj/cobra"
	flag "github.com/tj/pflag"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/utils"
)

// funcCommands is a list of commands which
// accept function names as arguments.
var funcCommands = map[string]bool{
	"build":    true,
	"delete":   true,
	"deploy":   true,
	"invoke":   true,
	"list":     true,
	"logs":     true,
	"metrics":  true,
	"rollback": true,
}

// Command config.
var Command = &cobra.Command{
	Use:              "autocomplete",
	Short:            "Generate bash auto-completion words",
	PersistentPreRun: root.PreRunNoop,
	Hidden:           true,
	Run:              run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// Run command.
func run(c *cobra.Command, args []string) {
	// TODO: works with flags inbetween?
	if len(args) == 0 {
		rootCommands()
		flags(root.Command)
		return
	}

	// find command, on error assume it's in-complete
	// or missing and all root commands are shown
	cmd := find(args[0])
	if cmd == nil {
		rootCommands()
		flags(root.Command)
		return
	}

	// command flags
	flags(cmd)

	// always show root flags
	flags(root.Command)

	// command accepts functions
	if ok := funcCommands[cmd.Name()]; ok {
		functions(args)
	}
}

// find command by `name`.
func find(name string) *cobra.Command {
	for _, cmd := range root.Command.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

// output root commands.
func rootCommands() {
	for _, cmd := range root.Command.Commands() {
		if !cmd.Hidden {
			fmt.Printf("%s ", cmd.Name())
		}
	}
}

// output flags.
func flags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if !f.Hidden {
			fmt.Printf("--%s ", f.Name)
		}
	})
}

// output (deduped) functions.
func functions(args []string) {
	files, err := ioutil.ReadDir("functions")
	if err != nil {
		return
	}

	for _, file := range files {
		if !utils.ContainsString(args, file.Name()) {
			fmt.Printf("%s ", file.Name())
		}
	}
}
