// Package logs outputs logs from CloudWatch logs.
package logs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/logs"
)

// filter pattern.
var filter string

// follow via polling.
var follow bool

// duration of results.
var duration time.Duration

// example output.
const example = `
    Print logs for all functions
    $ apex logs

    Follow the output
    $ apex logs -f

    Follow output with no historical logs
    $ apex logs -f --since 0

    Print logs for a single function
    $ apex logs api

    Print logs for all functions starting with "auth"
    $ apex logs auth*

    Print logs for functions with a specified start time, e.g. 5 minutes
    $ apex logs foo bar --since 5m`

// Command config.
var Command = &cobra.Command{
	Use:     "logs [<name>...] [<duration>]",
	Short:   "Output function logs",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.DurationVarP(&duration, "since", "s", 5*time.Minute, "Start time of the search")
	f.StringVarP(&filter, "filter", "F", "", "Filter logs with pattern")
	f.BoolVarP(&follow, "follow", "f", false, "Follow tails logs for updates")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	config := logs.Config{
		Service:       cloudwatchlogs.New(root.Session),
		StartTime:     time.Now().Add(-duration).UTC(),
		PollInterval:  5 * time.Second,
		Follow:        follow,
		FilterPattern: filter,
	}

	l := &logs.Logs{
		Config: config,
	}

	for _, fn := range root.Project.Functions {
		l.GroupNames = append(l.GroupNames, fn.GroupName())
	}

	for event := range l.Start() {
		fmt.Printf("\033[34m%s\033[0m %s", event.GroupName, event.Message)
	}

	return l.Err()
}
