package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/apex/log"
)

var buildCmdLocalValues struct {
	name string
}

const buildCmdExample = `  Build zip output for a function
  $ apex build foo > /tmp/out.zip`

var buildCmd = &cobra.Command{
	Use:     "build <name>",
	Short:   "Build a function",
	Example: buildCmdExample,
	PreRun:  buildCmdPreRun,
	Run:     buildCmdRun,
}

func buildCmdPreRun(c *cobra.Command, args []string) {
	lv := &buildCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}
	lv.name = args[0]
}

// outputs the generated archive to stdout
func buildCmdRun(c *cobra.Command, args []string) {
	lv := &buildCmdLocalValues

	err := pv.project.LoadFunctions(lv.name)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	zip, err := pv.project.Functions[0].Build()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	_, err = io.Copy(os.Stdout, zip)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
