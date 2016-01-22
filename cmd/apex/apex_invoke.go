package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/apex/apex/function"
	"github.com/apex/log"
)

type InvokeCmdLocalValues struct {
	Async bool
	Logs  bool

	name string
}

const invokeCmdExample = `  Invoke a function with input json
  $ apex invoke foo < request.json`

var invokeCmd = &cobra.Command{
	Use:     "invoke <name>",
	Short:   "Invoke a function",
	Example: invokeCmdExample,
	PreRun:  invokeCmdPreRun,
	Run:     invokeCmdRun,
}

var invokeCmdLocalValues = InvokeCmdLocalValues{}

func init() {
	lv := &invokeCmdLocalValues
	f := invokeCmd.Flags()

	f.BoolVarP(&lv.Async, "async", "a", false, "Async invocation")
	f.BoolVarP(&lv.Logs, "logs", "L", false, "Print logs")
}

func invokeCmdPreRun(c *cobra.Command, args []string) {
	lv := &invokeCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}
	lv.name = args[0]
}

// reads request json from stdin and outputs the responses
func invokeCmdRun(c *cobra.Command, args []string) {
	lv := &invokeCmdLocalValues
	dec := json.NewDecoder(os.Stdin)
	kind := function.RequestResponse

	if lv.Async {
		kind = function.Event
	}

	fn, err := pv.project.FunctionByName(lv.name)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	for {
		var v map[string]interface{}
		err := dec.Decode(&v)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error parsing response: %s", err)
		}

		var reply, logs io.Reader

		if e, ok := v["event"].(map[string]interface{}); ok {
			reply, logs, err = fn.Invoke(e, v["context"], kind)
		} else {
			reply, logs, err = fn.Invoke(v, nil, kind)
		}

		if err != nil {
			log.Fatalf("error response: %s", err)
		}

		if lv.Logs {
			io.Copy(os.Stderr, logs)
		}

		io.Copy(os.Stdout, reply)
		fmt.Fprintf(os.Stdout, "\n")
	}
}
