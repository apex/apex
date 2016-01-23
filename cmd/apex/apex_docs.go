package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/apex/wiki"
	"github.com/apex/log"
)

type DocsCmdLocalValues struct {
	topic string
}

const docsCmdExample = `  Output documentation topics
  $ apex docs

  Output documentation for a topic
  $ apex docs project.json`

var docsCmd = &cobra.Command{
	Use:              "docs [<topic>]",
	Short:            "Output documentation",
	Example:          docsCmdExample,
	PersistentPreRun: pv.noopRun,
	PreRun:           docsCmdPreRun,
	Run:              docsCmdRun,
}

var docsCmdLocalValues = DocsCmdLocalValues{}

func docsCmdPreRun(c *cobra.Command, args []string) {
	lv := &docsCmdLocalValues

	if len(args) >= 1 {
		lv.topic = strings.Join(args, " ")
	}
}

func docsCmdRun(c *cobra.Command, args []string) {
	lv := &docsCmdLocalValues

	var err error

	if lv.topic != "" {
		err = wiki.Topic(lv.topic, os.Stdout)
	} else {
		err = wiki.Topics(os.Stdout)
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
