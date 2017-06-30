// Package build outputs a function's zip to stdout.
package build

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/utils"
)

// name of function.
var name string

// env file.
var envFile string

// env supplied.
var env []string

// example output.
const example = `
    Build zip output for a function
    $ apex build foo > /tmp/out.zip`

// Command config.
var Command = &cobra.Command{
	Use:     "build <name>",
	Short:   "Build a function",
	Example: example,
	PreRunE: preRun,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.StringVarP(&envFile, "env-file", "E", "", "Set environment variables from JSON file")
	f.StringSliceVarP(&env, "set", "s", nil, "Set environment variable")
}

// PreRun errors if argument is missing.
func preRun(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing name argument")
	}

	name = args[0]
	return nil
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(name); err != nil {
		return err
	}

	fn := root.Project.Functions[0]

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

	zip, err := fn.Build()
	if err != nil {
		return err
	}

	if err := fn.Clean(); err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, zip)
	return err
}
