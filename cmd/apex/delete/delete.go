// Package delete resmove functions from Lambda.
package delete

import (
	"fmt"

	"github.com/tj/cobra"
	"github.com/tj/go-prompt"

	"github.com/apex/apex/cmd/apex/root"
)

// Force deletion.
var force bool

// example output.
const example = `
    Delete all functions
    $ apex delete

    Delete specified functions
    $ apex delete foo bar

    Delete all functions starting with "auth"
    $ apex delete auth*`

// Command config.
var Command = &cobra.Command{
	Use:     "delete [<name>...]",
	Short:   "Delete functions",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.BoolVarP(&force, "force", "f", false, "Force deletion")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	if !force && len(root.Project.Functions) > 1 {
		fmt.Printf("The following will be deleted:\n\n")
		for _, fn := range root.Project.Functions {
			fmt.Printf("  - %s\n", fn.Name)
		}
		fmt.Printf("\n")
	}

	if !force && !prompt.Confirm("Are you sure? (yes/no) ") {
		return nil
	}

	return root.Project.Delete()
}
