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

	err = global.ExecuteTemplate(&o, "page", baseContext(func(m map[string]any) {
		m["page"] = page
		m["title"] = page.Title
		m["guide"] = guide
		m["content"] = content
	}))
	if err != nil {
		log.Fatal("rendering page template", src, err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
}
