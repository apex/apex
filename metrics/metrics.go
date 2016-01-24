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

// getCloudWatchMetrics starts the CloudWatch metrics collection.
func getCloudWatchMetrics(fn string, startDate time.Time, endDate time.Time, service cloudwatchiface.CloudWatchAPI) {
	mc := &MetricCollector{
		Metrics:      []string{},
		Collected:    0,
		FunctionName: fn,
		Service:      service,
		StartDate:    startDate,
		EndDate:      endDate,
	}
	mc.Collect()
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
		go func(metricName string) {
			defer wg.Done()

			m := &Metric{
				metricName,
				nil,
			}

			var duration int
			switch x := mc.EndDate.Sub(mc.StartDate).Hours(); {
			default:
				duration = 1 // hourly
			case x > 24:
				duration = 24 // daily
			}

			period := time.Duration(duration) * time.Hour

			params := &cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(mc.EndDate),
				MetricName: aws.String(m.Name),
				Namespace:  aws.String("AWS/Lambda"),
				Period:     aws.Int64(int64(period.Seconds())),
				StartTime:  aws.Time(mc.StartDate),
				Statistics: []*string{
					aws.String("Sum"),
				},
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("FunctionName"),
						Value: aws.String(mc.FunctionName),
					},
				},
				Unit: aws.String(unit(metricName)),
			}

			resp, err := mc.Service.GetMetricStatistics(params)

			if err != nil {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
				return
			}

			m.Value = resp.Datapoints

			out <- *m
		}(name)
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

// unit for metric name.
func unit(name string) string {
	switch name {
	case "Duration":
		return "Milliseconds"
	default:
		return "Count"
	}
}
