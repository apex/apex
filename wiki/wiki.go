// Package wiki implements a simple GitHub wiki miner and output formatter.
package wiki

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/net/html"

	"github.com/apex/apex/colors"
)

// TODO: handle invalid page
// TODO: WikiTopics could be interactive

// Endpoint used to lookup wiki information.
var Endpoint = "https://github.com/apex/apex/wiki"

// Topics outputs topic categories from the wiki index page.
func Topics(w io.Writer) error {
	doc, err := goquery.NewDocument(Endpoint)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n")
	defer fmt.Fprintf(w, "\n")

	doc.Find(`#wiki-content .markdown-body ul li`).Each(func(i int, s *goquery.Selection) {
		strs := strings.Split(text(s), ": ")
		fmt.Fprintf(w, "  \033[%dm%s\033[0m: %s \n", colors.Blue, strs[0], strs[1])
	})

	fmt.Fprintf(w, "\n  Use `apex docs <topic>` to view a topic.\n")

	return nil
}

// Topic outputs documentation for the given `topic`'s page.
func Topic(topic string, w io.Writer) error {
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
		return fmt.Sprintf(" \033[%dm# %s\033[0m\n\n", colors.Blue, text(s))
	case node.Data == "h2":
		return fmt.Sprintf(" \033[%dm## %s\033[0m\n\n", colors.Blue, text(s))
	case node.Data == "h3":
		return fmt.Sprintf(" \033[%dm### %s\033[0m\n\n", colors.Blue, text(s))
	case node.Data == "p":
		return fmt.Sprintf("\033[%dm%s\033[0m\n\n", colors.None, indent(text(s), 1))
	case node.Data == "pre" || s.HasClass("highlight"):
		return fmt.Sprintf("\033[1m%s\033[0m\n\n", indent(text(s), 2))
	case node.Data == "a":
		return fmt.Sprintf("%s (%s) ", s.Text(), s.AttrOr("href", "missing link"))
	case node.Data == "li":
		return fmt.Sprintf("  • %s\n", contents(s))
	case node.Data == "ul":
		return fmt.Sprintf("%s\n", nodes(s))
	case node.Data == "code":
		return fmt.Sprintf("\033[1m%s\033[0m ", s.Text())
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
