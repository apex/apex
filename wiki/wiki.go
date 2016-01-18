// Package wiki implements a simple GitHub wiki miner and output formatter.
package wiki

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/net/html"
)

// TODO: handle invalid page
// TODO: WikiTopics could be interactive
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

// Endpoint used to lookup wiki information.
var Endpoint = "https://github.com/apex/apex/wiki"

// WikiTopics outputs topic categories.
func WikiTopics(w io.Writer) error {
	doc, err := goquery.NewDocument(Endpoint)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n")
	defer fmt.Fprintf(w, "\n")

	doc.Find(`#wiki-content .markdown-body ul li`).Each(func(i int, s *goquery.Selection) {
		strs := strings.Split(text(s), ": ")
		fmt.Fprintf(w, "  \033[%dm%s\033[0m: %s \n", blue, strs[0], strs[1])
	})

	fmt.Fprintf(w, "\n  Use `apex wiki <topic>` to view a topic.\n")

	return nil
}

// WikiTopic outputs topic for the given `topic`.
func WikiTopic(topic string, w io.Writer) error {
	doc, err := goquery.NewDocument(fmt.Sprintf("%s/%s", Endpoint, topic))
	if err != nil {
		return err
	}

	fmt.Fprintln(w)
	fmt.Fprint(w, nodes(doc.Find(`#wiki-content .markdown-body`)))
	fmt.Fprintln(w)

	return nil
}

// nodes returns a string representation of the selection's children.
func nodes(s *goquery.Selection) string {
	return strings.Join(s.Children().Map(node), "")
}

// contents returns a string representation of the selection's contents.
func contents(s *goquery.Selection) string {
	return strings.Join(s.Contents().Map(node), "")
}

// node returns a string representation of the selection.
func node(i int, s *goquery.Selection) string {
	switch node := s.Get(0); {
	case node.Data == "h1":
		return fmt.Sprintf(" \033[%dm# %s\033[0m\n\n", blue, text(s))
	case node.Data == "h2":
		return fmt.Sprintf(" \033[%dm## %s\033[0m\n\n", blue, text(s))
	case node.Data == "h3":
		return fmt.Sprintf(" \033[%dm### %s\033[0m\n\n", blue, text(s))
	case node.Data == "p":
		return fmt.Sprintf("\033[%dm%s\033[0m\n\n", none, indent(text(s), 1))
	case node.Data == "pre" || s.HasClass("highlight"):
		return fmt.Sprintf("\033[%dm%s\033[0m\n\n", gray, indent(text(s), 2))
	case node.Data == "a":
		return fmt.Sprintf("%s (%s) ", s.Text(), s.AttrOr("href", "missing link"))
	case node.Data == "li":
		return fmt.Sprintf("  • %s\n", contents(s))
	case node.Data == "ul":
		return fmt.Sprintf("%s\n", nodes(s))
	case node.Type == html.TextNode:
		return strings.TrimSpace(node.Data)
	default:
		return ""
	}
}

// text of selection, trimmed and wrapped.
func text(s *goquery.Selection) string {
	return wordwrap.WrapString(strings.TrimSpace(s.Text()), 80)
}

// indent string N times.
func indent(s string, n int) string {
	i := strings.Repeat("  ", n)
	return i + strings.Replace(s, "\n", "\n"+i, -1)
}
