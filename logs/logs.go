// Package logs implements a Kinesis log tailer.
package logs

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

// Tailer is a Kinesis log tailer.
type Tailer struct {
	Stream       string                  // Stream to tail.
	Service      kinesisiface.KinesisAPI // Service implementation.
	PollInterval time.Duration           // PollInterval [1s].

	mu  sync.Mutex
	err error
}

// Err returns the error (if any).
func (t *Tailer) Err() error {
	return t.err
}

// Start tailing and return channel of records.
func (t *Tailer) Start() <-chan *kinesis.Record {
	if t.PollInterval == 0 {
		t.PollInterval = time.Second
	}

	ch := make(chan *kinesis.Record)
	go t.start(ch)
	return ch
}

// start tailing by fetching the shards.
func (t *Tailer) start(ch chan<- *kinesis.Record) {
	res, err := t.Service.DescribeStream(&kinesis.DescribeStreamInput{
		StreamName: &t.Stream,
	})

	if err != nil {
		t.fail(err, ch)
		return
	}

	for _, shard := range res.StreamDescription.Shards {
		go t.consume(shard, ch)
	}
}

// consume a single shard.
func (t *Tailer) consume(shard *kinesis.Shard, ch chan<- *kinesis.Record) {
	res, err := t.Service.GetShardIterator(&kinesis.GetShardIteratorInput{
		ShardId:           shard.ShardId,
		StreamName:        &t.Stream,
		ShardIteratorType: aws.String("LATEST"),
	})

	if err != nil {
		t.fail(err, ch)
		return
	}

	iter := *res.ShardIterator

	for {
		res, err := t.Service.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: aws.String(iter),
		})

		if err != nil {
			t.fail(err, ch)
			return
		}

		iter = *res.NextShardIterator

		for _, record := range res.Records {
			ch <- record
		}

		time.Sleep(t.PollInterval)
	}
}

// fail with `err`.
func (t *Tailer) fail(err error, ch chan<- *kinesis.Record) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.err != nil {
		return
	}

	t.err = err
	close(ch)
}
