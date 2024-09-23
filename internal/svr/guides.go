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
	m.HandleFunc("/guides", renderGuide())
}

func renderGuide() http.HandlerFunc {
	var o bytes.Buffer
	err := global.ExecuteTemplate(&o, "guide", baseContext(func(m map[string]any) {
		m["guides"] = guides
		m["title"] = "Guides | Learn vince by example"
		m["meta"] = []Meta{
			{"description", "Useful tips and samples for working with vince"},
			{"og:site_name", "Vince"},
			{"og:title", "Guides | Learn vince by example"},
			{"og:url", "https://vinceanalytics.com/guides"},
			{"og:image", "https://vinceanalytics.com/images/logo.png"},
			{"og:locale", "en_US"},
			{"og:type", "website"},
			{"twitter:site", "@gernesti"},
			{"twitter:card", "summary_large_image"},
		}
	}))
	if err != nil {
		log.Fatal("rendering guide template", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
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
		m["meta"] = []Meta{
			{"og:site_name", "Vince"},
			{"og:title", page.Title},
			{"og:url", "https://vinceanalytics.com" + page.Link},
			{"og:image", "https://vinceanalytics.com/images/logo.png"},
			{"og:locale", "en_US"},
			{"og:type", "website"},
			{"twitter:site", "@gernesti"},
			{"twitter:card", "summary_large_image"},
		}
	}))
	if err != nil {
		log.Fatal("rendering page template", src, err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
}
