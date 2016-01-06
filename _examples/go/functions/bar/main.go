package main

import (
	"encoding/json"

	"github.com/apex/apex"
)

type Message struct {
	Hello string `json:"hello"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		return &Message{"bar"}, nil
	})
}
