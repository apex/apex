package apex_test

import (
	"encoding/json"

	"github.com/apex/apex"
)

type Message struct {
	Hello string `json:"hello"`
}

func Example() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		return &Message{"world"}, nil
	})
}
