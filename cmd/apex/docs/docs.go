// Package docs outputs Wiki documentation.
package docs

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	doc "github.com/apex/apex/docs"
)

// topic name.
var topic string

// example output.
const example = `  Output documentation topics
  $ apex docs

  Output documentation for a topic
  $ apex docs project.json`

// Command config.
var Command = &cobra.Command{
	Use:              "docs [<topic>]",
	Short:            "Output documentation",
	Example:          example,
	PersistentPreRun: root.PreRunNoop,
	PreRun:           preRun,
	RunE:             run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// PreRun joins args to form the topic.
func preRun(c *cobra.Command, args []string) {
	if len(args) >= 1 {
		topic = strings.Join(args, " ")
	}
}

// Run command.
func run(c *cobra.Command, args []string) (err error) {
	var w io.WriteCloser = os.Stdout

	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd := exec.Command("less", "-R")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		defer cmd.Wait()

		w, err = cmd.StdinPipe()
		if err != nil {
			return err
		}
		defer w.Close()

		if err := cmd.Start(); err != nil {
			return err
		}
	}

	_, err = io.Copy(w, doc.Reader())
	return err
}
