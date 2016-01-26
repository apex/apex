package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const listCmdExample = `  List all functions
  $ apex list`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List functions",
	Example: listCmdExample,
	Run:     listCmdRun,
}

func listCmdRun(c *cobra.Command, args []string) {
	fmt.Println()

	err := pv.project.LoadFunctions()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	for _, fn := range pv.project.Functions {
		fmt.Printf("  - %s (%s)\n", fn.Name, fn.Runtime)
	}
	fmt.Println()
}
