// Package help implements a simple GitHub wiki miner and output formatter.
package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/go-wordwrap"
)

// TODO: HelpTopics could be interactive
// TODO: walk nodes to highlight inline code etc
// TODO: ~/.apex.json user config, use here for color mapping etc

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Endpoint used to lookup help information.
var Endpoint = "https://github.com/apex/apex/wiki"

// Help outputs topic categories.
func HelpTopics(w io.Writer) error {
	doc, err := goquery.NewDocument(Endpoint)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n")
	defer fmt.Fprintf(w, "\n")

	doc.Find(`#wiki-content .markdown-body ul li`).Each(func(i int, s *goquery.Selection) {
		strs := strings.Split(text(s.Text()), ": ")
		fmt.Fprintf(w, "  \033[%dm%s\033[0m: %s \n", blue, strs[0], strs[1])
	})

	fmt.Fprintf(w, "\n  Use `apex help <topic>` to view a topic.\n")

	return nil
}

// HelpTopic outputs topic for the given `topic`.
func HelpTopic(topic string, w io.Writer) error {
	doc, err := goquery.NewDocument(fmt.Sprintf("%s/%s", Endpoint, topic))
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n")
	defer fmt.Fprintf(w, "\n")

	doc.Find(`#wiki-content .markdown-body *`).Each(func(i int, s *goquery.Selection) {
		switch node := s.Get(0); node.Data {
		case "h1":
			fmt.Printf(" \033[%dm# %s\033[0m\n\n", blue, text(s.Text()))
		case "h2":
			fmt.Printf(" \033[%dm## %s\033[0m\n\n", blue, text(s.Text()))
		case "h3":
			fmt.Printf(" \033[%dm### %s\033[0m\n\n", blue, text(s.Text()))
		case "p":
			fmt.Printf("\033[%dm%s\033[0m\n\n", none, indent(text(s.Text()), 1))
		case "div", "pre":
			if s.HasClass("highlight") || node.Data == "pre" {
				fmt.Printf("\033[%dm%s\033[0m\n\n", gray, indent(text(s.Text()), 2))
			}
		}
	})

	return nil
}

// text trim and wrap.
func text(s string) string {
	return wordwrap.WrapString(strings.TrimSpace(s), 80)
}

// indent string N times.
func indent(s string, n int) string {
	i := strings.Repeat("  ", n)
	return i + strings.Replace(s, "\n", "\n"+i, -1)
}
