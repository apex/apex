// Package metrics outputs metrics for a function.
package metrics

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/metrics"
)

// name of function.
var name string

// duration of results.
var duration time.Duration

// aggregatedMetric is an arregated metric :).
type aggregatedMetric struct {
	Name  string
	Count int
}

// example output.
const example = `  Output the last day of metrics for a function
  $ apex metrics foo

  Output the last three days of metrics for a function
  $ apex metrics foo --duration 72h`

// Command config.
var Command = &cobra.Command{
	Use:     "metrics <name> [<duration>]",
	Short:   "Output function metrics",
	Example: example,
	PreRunE: preRun,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.DurationVarP(&duration, "duration", "d", 24*time.Hour, "Duration of metrics results")
}

// PreRun errors if the name is missing.
func preRun(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing name argument")
	}

	name = args[0]
	return nil
}

// aggregate accumulates the datapoints.
func aggregate(values []*cloudwatch.Datapoint) int {
	sum := 0.0

	for _, dp := range values {
		sum += *dp.Sum
	}

	return int(sum)
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Project.LoadFunctions(name); err != nil {
		return err
	}

	fn := root.Project.Functions[0]

	start := time.Now().UTC().Add(-duration)
	end := time.Now().UTC()

	mc := &metrics.MetricCollector{
		Metrics:      []string{"Invocations", "Errors", "Duration", "Throttles"},
		Service:      cloudwatch.New(root.Session),
		FunctionName: fn.FunctionName,
		StartDate:    start,
		EndDate:      end,
	}

	aggregated := make(map[string]aggregatedMetric)

	for n := range mc.Collect() {
		aggregated[n.Name] = aggregatedMetric{n.Name, aggregate(n.Value)}
	}

	println()
	defer println()

	for _, m := range aggregated {
		switch m.Name {
		case "Duration":
			fmt.Printf("  \033[34m%11s:\033[0m %vms\n", m.Name, m.Count)
		default:
			fmt.Printf("  \033[34m%11s:\033[0m %v\n", m.Name, m.Count)
		}
	}

	return nil
}
