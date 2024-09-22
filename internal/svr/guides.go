package svr

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/vinceanalytics/sitetools/data"
)

var guides []*Guide

type Guide struct {
	Section string  `json:"section"`
	Pages   []*Page `json:"pages"`
}

type Page struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Link   string `json:"link"`
}

func registerGuides(m *http.ServeMux) {
	guidesPath := filepath.Join(*root, "guides.json")
	data, err := os.ReadFile(guidesPath)
	if err != nil {
		log.Fatal("reading guide data", err)
	}
	err = json.Unmarshal(data, &guides)
	if err != nil {
		log.Fatal("decoding guide data", err)
	}
	for _, g := range guides {
		for _, p := range g.Pages {
			m.HandleFunc(p.Link, renderPage(g, p))
		}
	}
}

func renderPage(guide *Guide, page *Page) http.HandlerFunc {
	src := filepath.Join(*root, page.Source)
	md, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(src, err)
	}
	content := template.HTML(markdown.ToHTML(md, nil, nil))
	var o bytes.Buffer
	tpl, err := template.ParseFS(data.Templates, "templates/page.html")
	if err != nil {
		log.Fatal("parsing page template", src, err)
	}
	err = tpl.Execute(&o, map[string]any{
		"page":    page,
		"guide":   guide,
		"content": content,
		"footer":  footer(),
	})
	if err != nil {
		log.Fatal("rendering page template", src, err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
}

func guideIndex() template.HTML {
	tpl, err := template.ParseFS(data.Templates, "templates/guides.html")
	if err != nil {
		log.Fatal("parsing guides template", err)
	}
	var o bytes.Buffer
	err = tpl.Execute(&o, guides)
	if err != nil {
		log.Fatal("rendering guides template", err)
	}
	return template.HTML(o.String())
}
