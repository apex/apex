package alias

import (
	"errors"

	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/stats"
)

// alias name.
var alias string

// version name
var version string

// example output.
const example = `
    Alias prod as current version
    $ apex alias prod

    Alias prod as version
    $ apex alias -v version prod

    Alias prod as current version for specific function
    $ apex alias prod function

    Alias prod as version for specific function
    $ apex alias -v version prod function
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
	f.StringVarP(&version, "version", "v", "$LATEST", "Function version")
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
	stats.Track("Alias", map[string]interface{}{
		"has_version": version != "",
		"args":        len(args),
	})

	if err := root.Project.LoadFunctions(args[1:]...); err != nil {
		return err
	}

	return root.Project.CreateOrUpdateAlias(alias, version)
}
