package build

import (
	"bytes"
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/google/go-github/v63/github"
)

type Page struct {
	Layout      string        `json:"layout,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Excerpt     template.HTML `json:"excerpt,omitempty"`
	Permalink   string        `json:"permalink,omitempty"`
	Date        Date          `json:"date,omitempty"`
	Author      struct {
		Name    string `json:"name,omitempty"`
		Twitter string `json:"twitter,omitempty"`
	} `json:"author,omitempty"`
	Next     *Page                       `json:"-"`
	Previous *Page                       `json:"-"`
	Index    int                         `json:"index,omitempty"`
	Source   string                      `json:"-"`
	Content  template.HTML               `json:"-"`
	Releases []*github.RepositoryRelease `json:"-"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ts, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	d.Time = ts
	return nil
}

func (d *Date) Timestamp() string {
	return d.String()
}

func (d *Date) Nice() string {
	return d.Format("Jan 02, 2006")
}

func (p *Page) Render(ctx Context, out string, tpl *template.Template) {
	base := filepath.Join(out, p.Permalink)
	os.MkdirAll(base, 0755)
	b := new(bytes.Buffer)
	ctx["page"] = p
	fail(tpl.Execute(b, ctx), "rendering page", "source", p.Source)
	dest := filepath.Join(base, "index.html")
	fmt.Println("=>", dest)
	fail(os.WriteFile(dest, b.Bytes(), 0600), "writing output file", "dest", dest)

}

func (p *Page) Read() {
	data, err := os.ReadFile(p.Source)
	fail(err, "reading page", "source", p.Source)
	md, err := jsonPre(p, data)
	fail(err, "reading front matter", "source", p.Source)
	p.Content = template.HTML(markdown.ToHTML(md, nil, nil))

	// set excerpt
	ex, _, _ := bytes.Cut(md, []byte("\n\n"))
	p.Excerpt = template.HTML(markdown.ToHTML(ex, nil, nil))
}

func (p *Page) URL() string {
	return p.Permalink
}

func Compare(a, b *Page) int {
	if !a.Date.IsZero() || !b.Date.IsZero() {
		// sort dated pages in descending order
		return b.Date.Compare(a.Date.Time)
	}
	// normally sort in the order they appeared in the file system
	return cmp.Compare(a.Index, b.Index)
}

type Pages []*Page

func (p Pages) Render(m Context, out string, tpl *template.Template) {
	for i := range p {
		p[i].Render(m, out, tpl)
	}
}

func (p Pages) Read() {
	for i := range p {
		p[i].Read()
	}
	slices.SortFunc(p, Compare)
}

type LayoutData map[string]Pages

func (layout LayoutData) Render(out string, tpl map[string]*template.Template) {
	m := make(Context)
	for k, v := range layout {
		m[k] = v
	}
	for k, v := range layout {
		v.Render(m, out, tpl[k])
	}
}

func (layout LayoutData) RenderSpecial(out string, pages []*Page, tpl map[string]*template.Template) {
	m := make(Context)
	for k, v := range layout {
		m[k] = v
	}
	for i := range pages {
		pages[i].Render(m, out, tpl[pages[i].Layout])
	}
}

type Context map[string]any

func (ctx Context) AbsoluteURL(args ...string) string {
	return absURL(args...)
}

var root = flag.String("url", "", "")

func absURL(args ...string) string {
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	a := strings.Join(args, "/")
	return *root + a
}

func Build(path string) LayoutData {
	pages, err := load(path)
	fail(err, "loading path", "path", path)
	layouts := make(LayoutData)
	for i := range pages {
		pages[i].Read()
		layouts[pages[i][0].Layout] = pages[i]
	}
	return layouts
}

func load(path string) ([]Pages, error) {
	var releases []*github.RepositoryRelease
	if data, err := os.ReadFile("releases.json"); err == nil {
		json.Unmarshal(data, &releases)
		slog.Info("found releases", "count", len(releases))
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	r := make([]Pages, 0, len(dir))
	for i := range dir {
		e := dir[i]
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(path, e.Name())

		child, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		if len(child) == 0 {
			continue
		}
		pages := make(Pages, 0, len(child))
		layout := e.Name()
		for j := range child {
			ch := child[j]
			if ch.IsDir() {
				continue
			}
			if filepath.Ext(ch.Name()) != ".md" {
				continue
			}

			full := filepath.Join(path, ch.Name())
			name := ch.Name()
			permalink := layout + "/" + name[:len(name)-3]

			pages = append(pages, &Page{
				Layout:    e.Name(),
				Index:     j,
				Source:    full,
				Permalink: permalink,
				Releases:  releases,
			})
		}
		if len(pages) == 0 {
			continue
		}
		r = append(r, pages)
	}
	return r, nil
}

func jsonPre(w any, b []byte) ([]byte, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	err := d.Decode(w)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(b[d.InputOffset():]), nil
}

func fail(err error, msg string, args ...any) {
	if err != nil {
		slog.Error(msg, append(args, slog.String("err", err.Error()))...)
		os.Exit(1)
	}
}

func SpecialPages() []*Page {
	return []*Page{
		{
			Layout:      "page",
			Title:       "Vince Analytics",
			Description: "Vince is a cloud native, self hosted, privacy-first alternative to google analytics",
			Permalink:   "/",
		},
		{
			Layout:      "post",
			Title:       "Vince Analytics Blog",
			Description: "Vince is a cloud native, self hosted, privacy-first alternative to google analytics",
			Permalink:   "/blog",
		},
	}
}
