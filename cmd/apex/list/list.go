// Package list outputs a list of Lambda function information.
package list

import (
	"fmt"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/colors"
	"github.com/apex/apex/stats"
)

// tfvars output format.
var tfvars bool

// example output.
const example = `
    List all functions
    $ apex list

    List functions based on glob
    $ apex list api_*

    Output list as Terraform variables (.tfvars)
    $ apex list --tfvars`

// Command config.
var Command = &cobra.Command{
	Use:     "list [<name>...]",
	Short:   "Output functions list",
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
	stats.Track("List", map[string]interface{}{
		"tfvars": tfvars,
	})

	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	if tfvars {
		outputTFvars()
	} else {
		outputList()
	}

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
		awsFn, err := fn.GetConfigCurrent()

		if awserr, ok := err.(awserr.Error); ok && awserr.Code() == "ResourceNotFoundException" {
			fmt.Printf("  \033[%dm%s\033[0m (not deployed) \n", colors.Blue, fn.Name)
		} else {
			fmt.Printf("  \033[%dm%s\033[0m\n", colors.Blue, fn.Name)
		}

		if fn.Description != "" {
			fmt.Printf("    description: %v\n", fn.Description)
		}
		fmt.Printf("    runtime: %v\n", fn.Runtime)
		fmt.Printf("    memory: %vmb\n", fn.Memory)
		fmt.Printf("    timeout: %vs\n", fn.Timeout)
		fmt.Printf("    role: %v\n", fn.Role)
		fmt.Printf("    handler: %v\n", fn.Handler)
		if awsFn != nil && awsFn.Configuration != nil && awsFn.Configuration.FunctionArn != nil {
			fmt.Printf("    arn: %v\n", *awsFn.Configuration.FunctionArn)
		}

		if err != nil {
			fmt.Println()
			continue // ignore
		}

		aliaslist, err := fn.GetAliases()
		if err != nil {
			continue
		}

		var aliases string
		for index, alias := range aliaslist.Aliases {
			if index > 0 {
				aliases += ", "
			}
			aliases += fmt.Sprintf("%s@v%s", *alias.Name, *alias.FunctionVersion)
		}
		if aliases == "" {
			aliases = "<none>"
		}
		fmt.Printf("    aliases: %s\n", aliases)
		fmt.Println()
	}
}
