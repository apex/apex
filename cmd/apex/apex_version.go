package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of Apex",
	Run:   versionCmdRun,
}

func versionCmdRun(c *cobra.Command, args []string) {
	fmt.Printf("Apex version %s\n", version)
}
