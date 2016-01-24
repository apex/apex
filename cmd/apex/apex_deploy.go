package main

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/log"
)

type DeployCmdLocalValues struct {
	concurrency int
	names       []string
}

const deployCmdExample = `  Deploy all functions
  $ apex deploy

  Deploy specific functions
  $ apex deploy foo bar

  Deploy functions in a different project
  $ apex deploy -C ~/dev/myapp`

var deployCmd = &cobra.Command{
	Use:     "deploy [<name>...]",
	Short:   "Deploy code and config changes",
	Example: deployCmdExample,
	PreRun:  deployCmdPreRun,
	Run:     deployCmdRun,
}

var deployCmdLocalValues = DeployCmdLocalValues{}

func init() {
	lv := &deployCmdLocalValues
	f := deployCmd.Flags()

	f.IntVarP(&lv.concurrency, "concurrency", "c", 5, "Concurrent deploys")
}

func deployCmdPreRun(c *cobra.Command, args []string) {
	lv := &deployCmdLocalValues

	if len(args) == 0 {
		lv.names = pv.project.FunctionNames()
		return
	}
	lv.names = args
}

func deployCmdRun(c *cobra.Command, args []string) {
	lv := &deployCmdLocalValues

	pv.project.Concurrency = lv.concurrency

	for _, s := range pv.Env {
		parts := strings.Split(s, "=")
		pv.project.Setenv(parts[0], parts[1])
	}

	if err := pv.project.DeployAndClean(lv.names); err != nil {
		log.Fatalf("error: %s", err)
	}
}
