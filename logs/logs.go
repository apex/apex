// Package logs implements AWS CloudWatchLogs tailing.
package logs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"

	"github.com/apex/apex/project"
	"github.com/apex/log"
)

// Logs implements log tailing for CloudWatchLogs.
type Logs struct {
	Service       cloudwatchlogsiface.CloudWatchLogsAPI
	Log           log.Interface
	Project       *project.Project
	FunctionName  string
	FilterPattern string

	err error
}

// Fetch log events.
func (l *Logs) Fetch() ([]*cloudwatchlogs.FilteredLogEvent, error) {
	group, err := l.logGroupName()
	if err != nil {
		return nil, err
	}

	start := time.Now().Add(-time.Minute).UnixNano() / int64(time.Millisecond)

	l.Log.Debugf("fetching %q with filter %q", group, l.FilterPattern)

	res, err := l.Service.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  &group,
		FilterPattern: &l.FilterPattern,
		StartTime:     &start,
	})

	if err != nil {
		return nil, err
	}

	return res.Events, nil
}

// Tail logs, make sure to check Err() after the returned channel closes.
func (l *Logs) Tail() <-chan *cloudwatchlogs.FilteredLogEvent {
	ch := make(chan *cloudwatchlogs.FilteredLogEvent)
	go l.loop(ch)
	return ch
}

// loop polls for log tailing.
func (l *Logs) loop(ch chan<- *cloudwatchlogs.FilteredLogEvent) {
	defer close(ch)

	var nextToken *string
	start := time.Now().Add(-time.Minute).UnixNano() / int64(time.Millisecond)

	group, err := l.logGroupName()
	if err != nil {
		l.err = err
		return
	}

	l.Log.Debugf("tailing %q with filter %q", group, l.FilterPattern)

	for {
		l.Log.Debugf("tailing from %d", start)

		var res *cloudwatchlogs.FilterLogEventsOutput
		var err error

		res, err = l.Service.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:  &group,
			FilterPattern: &l.FilterPattern,
			StartTime:     &start,
			NextToken:     nextToken,
		})

		if err != nil {
			l.err = err
			return
		}

		nextToken = res.NextToken

		for _, event := range res.Events {
			start = *event.Timestamp + 1
			ch <- event
		}

		time.Sleep(time.Second)
	}
}

// Err returns the first error, if any, during processing.
func (l *Logs) Err() error {
	return l.err
}

// logGroupName returns CloudWatch log group name.
func (l *Logs) logGroupName() (string, error) {
	fn, err := l.Project.FunctionByName(l.FunctionName)
	if err != nil {
		return "", err
	}
	fnName, err := l.Project.RenderFunctionName(fn)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/aws/lambda/%s", fnName), nil
}
