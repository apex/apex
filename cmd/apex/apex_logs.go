package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/apex/log"
)

type LogsCmdLocalValues struct {
	Filter   string
	Follow   bool
	Duration string

	name string
}

const logsCmdExample = `  Print logs for a function
  $ apex logs <name>

  Print logs for a function with a specified duration, e.g. 5 minutes
  $ apex logs <name> 5m`

var logsCmd = &cobra.Command{
	Use:     "logs <name> [<duration>]",
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
	f.BoolVarP(&lv.Follow, "follow", "f", false, "Tail logs")
}

func logsCmdPreRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}
	lv.name = args[0]

	if len(args) >= 2 {
		lv.Duration = args[1]
	}
}

func logsCmdRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues

	err := pv.project.LoadFunctions(lv.name)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	l, err := pv.project.Logs(pv.session, lv.Filter, lv.Duration)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if lv.Follow {
		for event := range l.Tail() {
			fmt.Printf("%s", *event.Message)
		}

		if err := l.Err(); err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	events, err := l.Fetch()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	for _, event := range events {
		fmt.Printf("%s", *event.Message)
	}
}
