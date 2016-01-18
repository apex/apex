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

	l := logs.Logs{
		Service:       service,
		Log:           log.Log,
		Project:       pv.project,
		FunctionName:  lv.name,
		FilterPattern: lv.Filter,
	}

	for event := range l.Tail() {
		fmt.Printf("%s", *event.Message)
	}

	if err := l.Err(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
