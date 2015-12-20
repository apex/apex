package main

import (
	"log"

	"github.com/apex/apex"
	"github.com/apex/apex/kinesis"
)

func main() {
	apex.Handle(kinesis.HandlerFunc(func(event *kinesis.Event, ctx *apex.Context) error {
		log.Printf("processing %d records", len(event.Records))

		for i, record := range event.Records {
			log.Printf("%d) %s", i, record.Data())
		}

		return nil
	}))
}
