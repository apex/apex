// Package deploy builds & deploys function zips.
package deploy

import (
	"fmt"

	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/utils"
)

// env vars.
var env []string

// env file.
var envFile string

// concurrency of deploys.
var concurrency int

// alias.
var alias string

// zip path.
var zip string

// example output.
const example = `
    Deploy all functions
    $ apex deploy

    Deploy specific functions
    $ apex deploy foo bar

    Deploy canary alias
    $ apex deploy foo --alias canary

    Deploy functions in a different project
    $ apex deploy -C ~/dev/myapp

    Deploy function with existing zip
    $ apex build > out.zip && apex deploy foo --zip out.zip

    Deploy all functions starting with "auth"
    $ apex deploy auth*`

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
	f.StringVarP(&envFile, "env-file", "E", "", "Set environment variables from JSON file")
	f.StringVarP(&alias, "alias", "a", "current", "Function alias")
	f.StringVarP(&zip, "zip", "z", "", "Zip path")
	f.IntVarP(&concurrency, "concurrency", "c", 5, "Concurrent deploys")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	root.Project.Concurrency = concurrency
	root.Project.Alias = alias
	root.Project.Zip = zip

	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	if envFile != "" {
		if err := root.Project.LoadEnvFromFile(envFile); err != nil {
			return fmt.Errorf("reading env file %q: %s", envFile, err)
		}
	}

	vars, err := utils.ParseEnv(env)
	if err != nil {
		return err
	}

	for k, v := range vars {
		root.Project.Setenv(k, v)
	}

	return root.Project.DeployAndClean()
}
