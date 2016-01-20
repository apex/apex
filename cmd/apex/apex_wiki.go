package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/apex/wiki"
	"github.com/apex/log"
)

type WikiCmdLocalValues struct {
	topic string
}

const wikiCmdExample = `  Output wiki topics
  $ apex wiki

  Output wiki for a topic
  $ apex wiki project.json`

var wikiCmd = &cobra.Command{
	Use:              "wiki [<topic>]",
	Short:            "Output wiki page pulled from the GitHub wiki",
	Example:          wikiCmdExample,
	PersistentPreRun: pv.noopRun,
	PreRun:           wikiCmdPreRun,
	Run:              wikiCmdRun,
}

var wikiCmdLocalValues = WikiCmdLocalValues{}

func wikiCmdPreRun(c *cobra.Command, args []string) {
	lv := &wikiCmdLocalValues

	if len(args) >= 1 {
		lv.topic = strings.Join(args, " ")
	}
}

func wikiCmdRun(c *cobra.Command, args []string) {
	lv := &wikiCmdLocalValues

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
