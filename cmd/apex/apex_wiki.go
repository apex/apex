package main

import (
	"os"

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
  # apex wiki project.json`

var wikiCmd = &cobra.Command{
	Use:     "wiki",
	Short:   "Output wiki page pulled from the GitHub wiki",
	Example: wikiCmdExample,
	PreRun:  wikiCmdPreRun,
	Run:     wikiCmdRun,
}

var wikiCmdLocalValues = WikiCmdLocalValues{}

func wikiCmdPreRun(c *cobra.Command, args []string) {
	lv := &wikiCmdLocalValues

	if len(args) >= 1 {
		lv.topic = args[0]
	}
}

func wikiCmdRun(c *cobra.Command, args []string) {
	lv := &wikiCmdLocalValues

	var err error

	if lv.topic != "" {
		err = wiki.WikiTopic(lv.topic, os.Stdout)
	} else {
		err = wiki.WikiTopics(os.Stdout)
	}

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
