// Package init bootstraps an Apex project.
package init

import (
	"github.com/spf13/cobra"

	"github.com/apex/apex/boot"
	"github.com/apex/apex/cmd/apex/root"
)

// example output.
const example = `  Initialize a project
  $ apex init`

// Command config.
var Command = &cobra.Command{
	Use:              "init",
	Short:            "Initialize a project",
	Example:          example,
	PersistentPreRun: root.PreRunNoop,
	RunE:             run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Prepare(c, args); err != nil {
		return err
	}

	return boot.All()
}
