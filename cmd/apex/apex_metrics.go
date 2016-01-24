package main

import (
	"fmt"
	"time"

	"github.com/apex/apex/metrics"
	"github.com/spf13/cobra"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type aggregatedMetric struct {
	Name  string
	Count int
}

type MetricsCmdLocalValues struct {
	Start    string
	End      string
	name     string
	duration string
}

const metricsCmdExample = `  Output the last day of metrics for a function
  $ apex metrics foo

  Output the last three days of metrics for a function
  $ apex metrics foo 72d`

var metricsCmd = &cobra.Command{
	Use:     "metrics <name> [<duration>]",
	Short:   "Output the CloudWatch metrics for a function",
	Example: metricsCmdExample,
	PreRun:  metricsCmdPreRun,
	Run:     metricsCmdRun,
}

var metricsCmdLocalValues = MetricsCmdLocalValues{}

func init() {
	lv := &metricsCmdLocalValues
	f := metricsCmd.Flags()

	f.StringVar(&lv.Start, "start", "", "Start Date")
	f.StringVar(&lv.End, "end", "", "End Date")
}

func metricsCmdPreRun(c *cobra.Command, args []string) {
	lv := &metricsCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}

	lv.name = args[0]

	if len(args) > 1 {
		lv.duration = args[1]
	} else {
		lv.duration = "24h"
	}
}

// aggregate accumulates the datapoints.
func aggregate(values []*cloudwatch.Datapoint) int {
	aggregated_sum := 0.0
	for _, dp := range values {
		aggregated_sum += *dp.Sum
	}
	return int(aggregated_sum)
}

func metricsCmdRun(c *cobra.Command, args []string) {
	lv := &metricsCmdLocalValues

	fn, err := pv.project.FunctionByName(lv.name)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	d, err := time.ParseDuration(lv.duration)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	start := time.Now().UTC().Add(-d)
	end := time.Now().UTC()

	mc := &metrics.MetricCollector{
		Metrics:      []string{"Invocations", "Errors", "Duration", "Throttles"},
		Collected:    0,
		FunctionName: fn.FunctionName,
		Service:      cloudwatch.New(pv.session),
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
}
