package metrics

import (
	"math/rand"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/golang/mock/gomock"

	"github.com/apex/apex/metrics/mock"
)

func TestGetStatistics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	serviceMock := mock_cloudwatchiface.NewMockCloudWatchAPI(mockCtrl)

	metrics := [][]string{
		[]string{"Invocations", "Count"},
		[]string{"Errors", "Count"},
		[]string{"Duration", "Milliseconds"},
		[]string{"Throttles", "Count"},
	}

	startTime := time.Date(2016, time.January, 17, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2016, time.January, 18, 18, 0, 0, 0, time.UTC)

	for _, metric := range metrics {
		rand.Seed(time.Now().UTC().UnixNano())
		serviceMock.EXPECT().GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
			MetricName: aws.String(metric[0]),
			Namespace:  aws.String("AWS/Lambda"),
			StartTime:  aws.Time(startTime),
			EndTime:    aws.Time(endTime),
			Period:     aws.Int64(int64((time.Duration(24) * time.Hour).Seconds())),
			Statistics: []*string{
				aws.String("Sum"),
			},
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String("FunctionName"),
					Value: aws.String("go_testf"),
				},
			},
			Unit: aws.String(metric[1]),
		}).Return(&cloudwatch.GetMetricStatisticsOutput{
			Datapoints: []*cloudwatch.Datapoint{
				&cloudwatch.Datapoint{Sum: aws.Float64(float64(rand.Intn(9999)))},
			},
			Label: aws.String("label"),
		}, nil)
	}

	name := "go_testf"

	mc := &MetricCollector{
		Metrics:      []string{"Invocations", "Errors", "Duration", "Throttles"},
		Collected:    0,
		FunctionName: name,
		Service:      serviceMock,
		StartDate:    startTime,
		EndDate:      endTime,
	}

	count := 0
	for _ = range mc.Collect() {
		count++
	}
	if count != 4 {
		t.Errorf("Wrong metrics count")
	}
}
