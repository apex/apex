package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tj/go-prompt"

	"github.com/apex/log"
)

var deleteCmdLocalValues struct {
	Force bool
}

const deleteCmdExample = `  Delete all functions
  $ apex delete

  Delete specified functions
  $ apex delete foo bar`

var deleteCmd = &cobra.Command{
	Use:     "delete [<name>...]",
	Short:   "Delete functions",
	Example: deleteCmdExample,
	Run:     deleteCmdRun,
}

func init() {
	lv := &deleteCmdLocalValues
	f := deleteCmd.Flags()

	f.BoolVarP(&lv.Force, "force", "f", false, "Force deletion")
}

func deleteCmdRun(c *cobra.Command, args []string) {
	lv := &deleteCmdLocalValues

	if err := pv.project.LoadFunctions(args...); err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	if !lv.Force && len(pv.project.Functions) > 1 {
		fmt.Printf("The following will be deleted:\n\n")
		for _, fn := range pv.project.Functions {
			fmt.Printf("  - %s\n", fn.Name)
		}
		fmt.Printf("\n")
	}

	if !lv.Force && !prompt.Confirm("Are you sure? (yes/no)") {
		return
	}

	if err := pv.project.Delete(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
