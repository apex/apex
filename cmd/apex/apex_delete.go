package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tj/go-prompt"

	"github.com/apex/log"
)

type DeleteCmdLocalValues struct {
	Force bool

	names []string
}

const deleteCmdExample = `  Delete all functions
  $ apex delete

  Delete specified functions
  $ apex delete foo bar`

var deleteCmd = &cobra.Command{
	Use:     "delete [<name>...]",
	Short:   "Delete functions",
	Example: deleteCmdExample,
	PreRun:  deleteCmdPreRun,
	Run:     deleteCmdRun,
}

var deleteCmdLocalValues = DeleteCmdLocalValues{}

func init() {
	lv := &deleteCmdLocalValues
	f := deleteCmd.Flags()

	f.BoolVarP(&lv.Force, "force", "f", false, "Force deletion")
}

func deleteCmdPreRun(c *cobra.Command, args []string) {
	lv := &deleteCmdLocalValues

	if len(args) == 0 {
		lv.names = pv.project.FunctionNames()
		return
	}
	lv.names = args
}

func deleteCmdRun(c *cobra.Command, args []string) {
	lv := &deleteCmdLocalValues

	if !lv.Force && len(lv.names) > 1 {
		fmt.Printf("The following will be deleted:\n\n")
		for _, name := range lv.names {
			fmt.Printf("  - %s\n", name)
		}
		fmt.Printf("\n")
	}

	if !lv.Force && !prompt.Confirm("Are you sure? (yes/no)") {
		return
	}

	if err := pv.project.Delete(lv.names); err != nil {
		log.Fatalf("error: %s", err)
	}
}
