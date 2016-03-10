// Package deploy builds & deploys function zips.
package deploy

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// env supplied.
var env []string

// concurrency of deploys.
var concurrency int

// alias.
var alias string

// example output.
const example = `  Deploy all functions
  $ apex deploy

  Deploy specific functions
  $ apex deploy foo bar

  Deploy canary alias
  $ apex deploy foo --alias canary

  Deploy functions in a different project
  $ apex deploy -C ~/dev/myapp`

// Command config.
var Command = &cobra.Command{
	Use:     "deploy [<name>...]",
	Short:   "Deploy functions and config",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.StringSliceVarP(&env, "set", "s", nil, "Set environment variable")
	f.StringVarP(&alias, "alias", "a", "current", "Function alias")
	f.IntVarP(&concurrency, "concurrency", "c", 5, "Concurrent deploys")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	root.Project.Concurrency = concurrency
	root.Project.Alias = alias

	c.Root().PersistentFlags().Lookup("name string")

	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	for _, s := range env {
		parts := strings.Split(s, "=")
		root.Project.Setenv(parts[0], parts[1])
	}

	return root.Project.DeployAndClean()
}
