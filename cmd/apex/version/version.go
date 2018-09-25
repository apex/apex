// Package version oututs the version.
package version

import (
	"fmt"

	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// Version of program.
const Version = "master"

// Command config.
var Command = &cobra.Command{
	Use:              "version",
	Short:            "Print version of Apex",
	PersistentPreRun: root.PreRunNoop,
	Run:              run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// Run command.
func run(c *cobra.Command, args []string) {
	fmt.Printf("Apex version %s\n", Version)
}
