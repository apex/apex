package main

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"

	"github.com/apex/apex/logs"
)

var logsCmdLocalValues struct {
	Filter   string
	Follow   bool
	Duration time.Duration
}

const logsCmdExample = `  Print logs for all functions
  $ apex logs

  Follow the output
  $ apex logs -f

  Print logs for a single function
  $ apex logs api

  Print logs for functions with a specified duration, e.g. 5 minutes
  $ apex logs foo bar --duration 5m`

var logsCmd = &cobra.Command{
	Use:     "logs [<name>...] [<duration>]",
	Short:   "Output logs with optional filter pattern",
	Example: logsCmdExample,
	Run:     logsCmdRun,
}

func init() {
	lv := &logsCmdLocalValues
	f := logsCmd.Flags()

	f.DurationVarP(&lv.Duration, "duration", "d", 5*time.Minute, "Duration of log search prior to now")
	f.StringVarP(&lv.Filter, "filter", "F", "", "Filter logs with pattern")
	f.BoolVarP(&lv.Follow, "follow", "f", false, "Follow tails logs for updates")
}

func logsCmdRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues

	if err := pv.project.LoadFunctions(args...); err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	config := logs.Config{
		Service:       cloudwatchlogs.New(pv.session),
		FilterPattern: lv.Filter,
		PollInterval:  2 * time.Second,
		StartTime:     time.Now().Add(-lv.Duration).UTC(),
		Follow:        lv.Follow,
	}

	l := &logs.Logs{
		Config: config,
	}

	for _, fn := range pv.project.Functions {
		l.GroupNames = append(l.GroupNames, fn.GroupName())
	}

	for event := range l.Start() {
		fmt.Printf("\033[34m%s\033[0m %s", event.GroupName, event.Message)
	}

	if err := l.Err(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
