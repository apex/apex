//go:generate go-bindata -pkg docs .

// Package docs outputs colored markdown for a given topic.
package docs

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/golang-commonmark/markdown"
	"github.com/mitchellh/go-wordwrap"

	"github.com/apex/apex/colors"
)

// page of documentation.
type page struct {
	Name string
	File string
}

// pages by order.
var pages = []page{
	{"Installation", "installation.md"},
	{"AWS credentials", "aws-credentials.md"},
	{"Getting started", "getting-started.md"},
	{"Structuring projects", "projects.md"},
	{"Structuring functions", "functions.md"},
	{"Deploying functions", "deploy.md"},
	{"Invoking functions", "invoke.md"},
	{"Listing functions", "list.md"},
	{"Deleting functions", "delete.md"},
	{"Building functions", "build.md"},
	{"Rolling back versions", "rollback.md"},
	{"Function hooks", "hooks.md"},
	{"Viewing log output", "logs.md"},
	{"Viewing metrics", "metrics.md"},
	{"Managing infrastructure", "infra.md"},
	{"Previewing with dry-run", "dryrun.md"},
	{"Environment variables", "env.md"},
	{"Omitting files", "ignore.md"},
	{"Understanding the shim", "shim.md"},
	{"Viewing documentation", "docs.md"},
	{"Upgrading Apex", "upgrade.md"},
	{"FAQ", "faq.md"},
	{"Links", "links.md"},
}

// Reader returns all documentation as a single page.
func Reader() io.Reader {
	var in bytes.Buffer

	for _, page := range pages {
		in.WriteString(fmt.Sprintf("\n# %s\n", page.Name))
		in.Write(MustAsset(page.File))
	}

	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	v := &renderer{}
	s := v.visit(md.Parse(in.Bytes()))
	return strings.NewReader(s)
}

// indent string N times.
func indent(s string, n int) string {
	i := strings.Repeat("  ", n)
	return i + strings.Replace(s, "\n", "\n"+i, -1)
}

// renderer for terminal output.
type renderer struct {
	inParagraph bool
	inList      bool
	inLink      string
}

// visit `tokens`.
func (r *renderer) visit(tokens []markdown.Token) (s string) {
	for _, t := range tokens {
		s += r.visitToken(t)
	}
	return
}

// vistToken `t`.
func (r *renderer) visitToken(t markdown.Token) string {
	switch t.(type) {
	case *markdown.ParagraphOpen:
		r.inParagraph = true
		return ""
	case *markdown.ParagraphClose:
		r.inParagraph = false
		return "\n"
	case *markdown.CodeBlock:
		return fmt.Sprintf("\n%s\n", indent(t.(*markdown.CodeBlock).Content, 2))
	case *markdown.Fence:
		return fmt.Sprintf("\n%s\n", indent(t.(*markdown.Fence).Content, 2))
	case *markdown.HeadingOpen:
		n := t.(*markdown.HeadingOpen).HLevel
		return fmt.Sprintf("\n  %s \033[%dm", strings.Repeat("#", n), colors.Blue)
	case *markdown.HeadingClose:
		return "\n\033[0m\n"
	case *markdown.StrongOpen:
		return "\033[1m"
	case *markdown.StrongClose:
		return "\033[0m"
	case *markdown.BulletListOpen:
		r.inList = true
		return "\n"
	case *markdown.BulletListClose:
		r.inList = false
		return "\n"
	case *markdown.ListItemOpen:
		return "  - "
	case *markdown.LinkOpen:
		r.inLink = t.(*markdown.LinkOpen).Href
		return ""
	case *markdown.CodeInline:
		s := t.(*markdown.CodeInline).Content
		return fmt.Sprintf("\033[%dm%s\033[0m", colors.Gray, s)
	case *markdown.Text:
		s := t.(*markdown.Text).Content

		if r.inLink != "" {
			s = fmt.Sprintf("%s (%s)", s, r.inLink)
			r.inLink = ""
		}

		return s
	case *markdown.Inline:
		s := r.visit(t.(*markdown.Inline).Children)

		if r.inParagraph && !r.inList {
			s = indent(wordwrap.WrapString(s, 75), 1)
		}

		return s
	default:
		return ""
	}
}
