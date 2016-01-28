package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/apex/log"
)

var invokeCmdLocalValues struct {
	Logs bool

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

func init() {
	lv := &invokeCmdLocalValues
	f := invokeCmd.Flags()

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

	err := pv.project.LoadFunctions(lv.name)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}

	fn := pv.project.Functions[0]

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
			reply, logs, err = fn.Invoke(e, v["context"])
		} else {
			reply, logs, err = fn.Invoke(v, nil)
		}

		if lv.Logs {
			io.Copy(os.Stderr, logs)
		}

		if err != nil {
			log.Fatalf("error response: %s", err)
		}

		io.Copy(os.Stdout, reply)
		fmt.Fprintf(os.Stdout, "\n")
	}
}
