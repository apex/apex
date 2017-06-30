package alias

import (
	"errors"

	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// alias name.
var alias string

// version name
var version string

// example output.
const example = `
    Alias all functions as "prod"
    $ apex alias prod

    Alias all "api_*" functions to "prod"
    $ apex alias prod api_*

    Alias all functions of version 5 to "prod"
    $ apex alias -v v5 prod

    Alias specific function to "stage"
    $ apex alias stage myfunction

    Alias specific function's version 10 to "stage"
    $ apex alias -v v10 stage myfunction
`

// Command config.
var Command = &cobra.Command{
	Use:     "alias [<name>...]",
	Short:   "Create or update alias on functions",
	Example: example,
	PreRunE: preRun,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.StringVarP(&version, "version", "v", "current", "Function version")
}

// PreRun errors if the alias argument is missing.
func preRun(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing alias argument")
	}

	alias = args[0]
	return nil
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(args[1:]...); err != nil {
		return err
	}

	return root.Project.CreateOrUpdateAlias(alias, version)
}
