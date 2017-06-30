// Package invoke calls a Lambda function.
package invoke

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/tj/cobra"

	"github.com/apex/apex/cmd/apex/root"
)

// alias.
var alias string

// includeLogs in output.
var includeLogs bool

// name of function.
var name string

// example output.
const example = `
    Invoke a function with input json
    $ apex invoke foo < request.json

    Invoke canary alias
    $ apex invoke foo < request.json --alias canary`

// Command config.
var Command = &cobra.Command{
	Use:     "invoke <name>",
	Short:   "Invoke functions",
	Example: example,
	PreRunE: preRun,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.Flags()
	f.BoolVarP(&includeLogs, "logs", "L", false, "Print logs")
	f.StringVarP(&alias, "alias", "a", "current", "Function alias")
}

// PreRun errors if the name argument is missing.
func preRun(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing name argument")
	}

	name = args[0]
	return nil
}

// Run command.
func run(c *cobra.Command, args []string) error {
	dec := json.NewDecoder(input())

	root.Project.Alias = alias

	if err := root.Project.LoadFunctions(name); err != nil {
		return err
	}

	fn := root.Project.Functions[0]

	for {
		var v map[string]interface{}
		err := dec.Decode(&v)

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("parsing response: %s", err)
		}

		var reply, logs io.Reader

		if e, ok := v["event"].(map[string]interface{}); ok {
			reply, logs, err = fn.Invoke(e, v["context"])
		} else {
			reply, logs, err = fn.Invoke(v, nil)
		}

		if includeLogs && logs != nil {
			io.Copy(os.Stderr, logs)
		}

		if err != nil {
			return fmt.Errorf("function response: %s", err)
		}

		io.Copy(os.Stdout, reply)
		fmt.Fprintf(os.Stdout, "\n")
	}

	return nil
}

// input from stdin or empty object by default.
func input() io.Reader {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return strings.NewReader("{}")
	}

	return os.Stdin
}
