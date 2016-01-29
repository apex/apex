// Package list outputs a list of Lambda function information.
package list

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// tfvars output format.
var tfvars bool

// example output.
const example = `  List all functions
  $ apex list

  Output list as Terraform variables (.tfvars)
  $ apex list --tfvars`

// Command config.
var Command = &cobra.Command{
	Use:     "list",
	Short:   "List functions",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.BoolVar(&tfvars, "tfvars", false, "Output as Terraform variables")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(); err != nil {
		return err
	}

	if tfvars {
		outputTFvars()
	}

	outputList()
	return nil
}

// outputTFvars format.
func outputTFvars() {
	for _, fn := range root.Project.Functions {
		config, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}

		fmt.Printf("apex_function_%s=%q\n", fn.Name, *config.Configuration.FunctionArn)
	}
}

// outputList format.
func outputList() {
	fmt.Println()
	for _, fn := range root.Project.Functions {
		fmt.Printf("  - %s (%s)\n", fn.Name, fn.Runtime)
	}
	fmt.Println()
}
