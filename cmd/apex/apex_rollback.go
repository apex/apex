package main

import (
	"github.com/spf13/cobra"

	"github.com/apex/log"
)

type RollbackCmdLocalValues struct {
	name    string
	version string
}

const rollbackCmdExample = `  Rollback a function to the previous version
  $ apex rollback foo

  Rollback a function to the specified version
  $ apex rollback bar 3`

var rollbackCmd = &cobra.Command{
	Use:     "rollback <name>",
	Short:   "Rollback a function with optional version",
	Example: rollbackCmdExample,
	PreRun:  rollbackCmdPreRun,
	Run:     rollbackCmdRun,
}

var rollbackCmdLocalValues = RollbackCmdLocalValues{}

func rollbackCmdPreRun(c *cobra.Command, args []string) {
	lv := &rollbackCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}
	lv.name = args[0]

	if len(args) >= 2 {
		lv.version = args[1]
	}
}

func rollbackCmdRun(c *cobra.Command, args []string) {
	lv := &rollbackCmdLocalValues

	err := pv.project.LoadFunctions(lv.name)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	fn := pv.project.Functions[0]

	if lv.version == "" {
		err = fn.Rollback()
	} else {
		err = fn.RollbackVersion(lv.version)
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
