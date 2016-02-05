// Package metrics outputs metrics for a function.
package metrics

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/colors"
	"github.com/apex/apex/metrics"
)

// name of function.
var name string

// duration of results.
var duration time.Duration

// example output.
const example = `  Print the last 24 hours of metrics for all functions
  $ apex metrics

  Print the last 24 hours of metrics for a function
  $ apex metrics foo

  Print metrics for a function with a specified start time, e.g. the last 3 days
  $ apex metrics foo --start 72h`

// Command config.
var Command = &cobra.Command{
	Use:     "metrics [<name>...] [<duration>]",
	Short:   "Output function metrics",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.DurationVarP(&duration, "start", "s", 24*time.Hour, "Start time of the results")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(args...); err != nil {
		return err
	}

	config := metrics.Config{
		Service:   cloudwatch.New(root.Session),
		StartDate: time.Now().UTC().Add(-duration),
		EndDate:   time.Now().UTC(),
	}

	m := metrics.Metrics{
		Config: config,
	}

	for _, fn := range root.Project.Functions {
		m.FunctionNames = append(m.FunctionNames, fn.FunctionName)
	}

	aggregated := m.Collect()

	fmt.Println()
	for _, fn := range root.Project.Functions {
		fnMetrics := aggregated[fn.FunctionName]

		fmt.Printf("  \033[%dm%s\033[0m\n", colors.Blue, fn.Name)
		fmt.Printf("    invocations: %v\n", fnMetrics.Invocations)
		fmt.Printf("    duration: %vms\n", fnMetrics.Duration)
		fmt.Printf("    throttles: %v\n", fnMetrics.Throttles)
		fmt.Printf("    error: %v\n", fnMetrics.Errors)
		fmt.Println()
	}

	return nil
}
