package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"

	"github.com/apex/apex/logs"
	"github.com/apex/log"
)

type LogsCmdLocalValues struct {
	Filter string

	name string
}

const logsCmdExample = `  Print logs for a function
  $ apex logs <name>`

var logsCmd = &cobra.Command{
	Use:     "logs <name>",
	Short:   "Output logs with optional filter pattern",
	Example: logsCmdExample,
	PreRun:  logsCmdPreRun,
	Run:     logsCmdRun,
}

var logsCmdLocalValues = LogsCmdLocalValues{}

func init() {
	lv := &logsCmdLocalValues
	f := logsCmd.Flags()

	f.StringVarP(&lv.Filter, "filter", "F", "", "Filter logs with pattern")
}

func logsCmdPreRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}
	lv.name = args[0]
}

func logsCmdRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues
	service := cloudwatchlogs.New(pv.session)

	// TODO(tj): refactor logs.Logs to take Project so this hack
	// can be removed, it'll also make multi-function tailing easier
	group := fmt.Sprintf("/aws/lambda/%s_%s", pv.project.Name, lv.name)

	l := logs.Logs{
		LogGroupName:  group,
		FilterPattern: lv.Filter,
		Service:       service,
		Log:           log.Log,
	}

	for event := range l.Tail() {
		fmt.Printf("%s", *event.Message)
	}

	if err := l.Err(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
