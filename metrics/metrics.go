//go:generate mockgen -destination mock/cloudwatchiface.go github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface CloudWatchAPI

// Package metrics fetches CloudWatch metrics for a function.
package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// MetricCollector a wrapper all metrics for a specific function.
type MetricCollector struct {
	Metrics      []string
	Collected    int
	FunctionName string
	Service      cloudwatchiface.CloudWatchAPI
	StartDate    time.Time
	EndDate      time.Time
}

// Metric represents a CloudWatch metric with a given name and value.
type Metric struct {
	Name  string
	Value []*cloudwatch.Datapoint
}

// Collect builds the collector pipeline
func (mc *MetricCollector) Collect() <-chan Metric {
	return mc.collect(mc.gen())
}

// collect starts a new cloudwatch session and requests the key metrics.
func (mc *MetricCollector) collect(in <-chan string) <-chan Metric {
	var wg sync.WaitGroup
	out := make(chan Metric)

	for name := range in {
		wg.Add(1)
		name := name

		go func() {
			defer wg.Done()

			m := &Metric{Name: name}

			params := &cloudwatch.GetMetricStatisticsInput{
				StartTime:  aws.Time(mc.StartDate),
				EndTime:    aws.Time(mc.EndDate),
				MetricName: aws.String(m.Name),
				Namespace:  aws.String("AWS/Lambda"),
				Period:     aws.Int64(int64(period(mc.StartDate, mc.EndDate).Seconds())),
				Statistics: []*string{
					aws.String("Sum"),
				},
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("FunctionName"),
						Value: aws.String(mc.FunctionName),
					},
				},
				Unit: aws.String(unit(name)),
			}

			resp, err := mc.Service.GetMetricStatistics(params)

			if err != nil {
				// TODO: refactor so that errors are reported in cmd
				fmt.Println(err.Error())
				return
			}

			m.Value = resp.Datapoints

			out <- *m
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// gen generates the key metric structs and returns a channel pipeline.
func (mc *MetricCollector) gen() <-chan string {
	out := make(chan string, len(mc.Metrics))
	for _, n := range mc.Metrics {
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
