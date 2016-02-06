// Package metrics fetches CloudWatch metrics for a function.
package metrics

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// Config is used to configure Metrics.
type Config struct {
	MetricNames []string
	Service     cloudwatchiface.CloudWatchAPI
	StartDate   time.Time
	EndDate     time.Time
}

// AggregatedMetrics represents aggregated metrics.
type AggregatedMetrics struct {
	Duration    int
	Errors      int
	Invocations int
	Throttles   int
}

// Metrics collects CloudWatch metrics for multiple functions
type Metrics struct {
	Config
	FunctionNames []string
}

// Collect and aggregate metrics for multiple functions.
func (m *Metrics) Collect() (a map[string]AggregatedMetrics) {
	a = make(map[string]AggregatedMetrics)

	for _, fnName := range m.FunctionNames {
		metric := Metric{
			Config:       m.Config,
			FunctionName: fnName,
		}

		a[fnName] = metric.Collect()
	}

	return
}
