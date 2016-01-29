// Package build outputs a function's zip to stdout.
package build

import (
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// name of function.
var name string

// example output.
const example = `  Build zip output for a function
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
