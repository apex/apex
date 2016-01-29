package main

import (
	"fmt"

	"github.com/apex/log"

	"github.com/spf13/cobra"
)

var listCmdLocalValues struct {
	Tfvars bool
}

const listCmdExample = `  List all functions
  $ apex list

  Output list as Terraform variables (.tfvars)
  $ apex list --tfvars`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List functions",
	Example: listCmdExample,
	Run:     listCmdRun,
}

func init() {
	lv := &listCmdLocalValues
	f := listCmd.Flags()

	f.BoolVar(&lv.Tfvars, "tfvars", false, "Output as Terraform variables")
}

func listCmdRun(c *cobra.Command, args []string) {
	lv := &listCmdLocalValues

	err := pv.project.LoadFunctions()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if lv.Tfvars {
		outputTfvars()
	} else {
		outputList()
	}
}

func outputTfvars() {
	for _, fn := range pv.project.Functions {
		fnConfig, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}
		fmt.Printf("apex_function_%s=\"%s\"\n", fn.Name, *fnConfig.Configuration.FunctionArn)
	}
}

func outputList() {
	fmt.Println()
	for _, fn := range pv.project.Functions {
		fmt.Printf("  - %s (%s)\n", fn.Name, fn.Runtime)
	}
	fmt.Println()
}
