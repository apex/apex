package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var metricsNames = []string{"Invocations", "Errors", "Duration", "Throttles"}

// Metric collects metrics for single function.
type Metric struct {
	Config
	FunctionName string
}

// Collect and aggregate metrics for on function.
func (m *Metric) Collect() (a AggregatedMetrics) {
	for n := range m.collect(m.gen()) {
		value := aggregate(n.Value)

		switch n.Name {
		case "Duration":
			a.Duration = value
		case "Errors":
			a.Errors = value
		case "Invocations":
			a.Invocations = value
		case "Throttles":
			a.Throttles = value
		}
	}

	return
}

// cloudWatchMetric represents a CloudWatch metric with a given name and value.
type cloudWatchMetric struct {
	Name  string
	Value []*cloudwatch.Datapoint
}

// stats for function `name`.
func (m *Metric) stats(name string) (*cloudwatch.GetMetricStatisticsOutput, error) {
	return m.Service.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		StartTime:  &m.StartDate,
		EndTime:    &m.EndDate,
		MetricName: &name,
		Namespace:  aws.String("AWS/Lambda"),
		Period:     aws.Int64(int64(period(m.StartDate, m.EndDate).Seconds())),
		Statistics: []*string{
			aws.String("Sum"),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("FunctionName"),
				Value: &m.FunctionName,
			},
		},
		Unit: aws.String(unit(name)),
	})
}

// collect starts a new cloudwatch session and requests the key metrics.
func (m *Metric) collect(in <-chan string) <-chan cloudWatchMetric {
	var wg sync.WaitGroup
	out := make(chan cloudWatchMetric)

	for name := range in {
		wg.Add(1)
		name := name

		go func() {
			defer wg.Done()

			res, err := m.stats(name)
			if err != nil {
				// TODO: refactor so that errors are reported in cmd
				fmt.Println(err.Error())
				return
			}

			out <- cloudWatchMetric{
				Name:  name,
				Value: res.Datapoints,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// gen generates the key metric structs and returns a channel pipeline.
func (m *Metric) gen() <-chan string {
	out := make(chan string, len(metricsNames))
	for _, n := range metricsNames {
		out <- n
	}
	close(out)
	return out
}

// period returns the resolution of metrics.
func period(start, end time.Time) time.Duration {
	switch n := end.Sub(start).Hours(); {
	case n > 24:
		return time.Hour * 24
	default:
		return time.Hour
	}
}

// unit for metric name.
func unit(name string) string {
	switch name {
	case "Duration":
		return "Milliseconds"
	default:
		return "Count"
	}
}

// aggregate accumulates the datapoints.
func aggregate(datapoints []*cloudwatch.Datapoint) int {
	sum := 0.0

	for _, datapoint := range datapoints {
		sum += *datapoint.Sum
	}

	return int(sum)
}
