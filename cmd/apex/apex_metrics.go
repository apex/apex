package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/apex/apex/metrics"
	"github.com/spf13/cobra"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const timeFormat = "02/01/2006 15:04"

type aggregatedMetric struct {
	Name  string
	Count int
}

type MetricsCmdLocalValues struct {
	Start string
	End   string
	name  string
}

const metricsCmdExample = `  Output the CloudWatch metrics for a function for the last 24 hours time range
  $ apex metrics foo

  Output the CloudWatch metrics for a function for a customized time range
  $ apex metrics foo --start "18/01/2016 10:00" --end "19/01/2016 22:00"`

var metricsCmd = &cobra.Command{
	Use:     "metrics <name> [--start <startDate>] [--end <endDate>]",
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

	var sd string
	var ed string

	if lv.Start == "" {
		sd = time.Now().AddDate(0, 0, -1).Format(timeFormat)
	} else {
		sd = lv.Start
	}

	if lv.End == "" {
		ed = time.Now().Format(timeFormat)
	} else {
		ed = lv.End
	}

	s, _ := time.Parse(timeFormat, sd)
	e, _ := time.Parse(timeFormat, ed)

	mc := &metrics.MetricCollector{
		Metrics:      []string{"Invocations", "Errors", "Duration", "Throttles"},
		Collected:    0,
		FunctionName: fn.FunctionName,
		Service:      cloudwatch.New(pv.session),
		StartDate:    s,
		EndDate:      e,
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 9, 0, '\t', 0)
	collMetrics := make(map[string]aggregatedMetric)

	for n := range mc.Collect() {
		collMetrics[n.Name] = aggregatedMetric{n.Name, aggregate(n.Value)}
	}

	for _, m := range mc.Metrics {
		mm := collMetrics[m]
		switch {
		case m == "Duration":
			fmt.Fprintf(w, "\033[%dm%s:\033[0m\t%vms\n", 37, mm.Name, mm.Count)
		default:
			fmt.Fprintf(w, "\033[%dm%s:\033[0m\t%v\n", 37, mm.Name, mm.Count)
		}
	}
	w.Flush()

}
