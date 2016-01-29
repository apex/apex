// Package docs outputs Wiki documentation.
package docs

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/wiki"
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
func run(c *cobra.Command, args []string) error {
	if topic == "" {
		return wiki.Topics(os.Stdout)
	}

	return wiki.Topic(topic, os.Stdout)
}
