// Package rollback implements function version rollback.
package rollback

import (
	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// alias.
var alias string

// version target.
var version string

// example output.
const example = `
    Rollback all functions to the previous version
    $ apex rollback

    Rollback canary alias for a function
    $ apex rollback foo --alias canary

    Rollback all functions starting with "auth"
    $ apex rollback auth*

    Rollback a function to the specified version
    $ apex rollback bar -v 3`

// Command config.
var Command = &cobra.Command{
	Use:     "rollback [<name>...]",
	Short:   "Rollback functions",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.StringVarP(&alias, "alias", "a", "current", "Function alias")
	f.StringVarP(&version, "version", "v", "", "version to which rollback is done")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	root.Project.Alias = alias

	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	if version == "" {
		return root.Project.Rollback()
	}

	return root.Project.RollbackVersion(version)
}
