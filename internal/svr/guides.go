package svr

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gomarkdown/markdown"
)

var guides []*Guide

type Guide struct {
	Section string  `json:"section"`
	Pages   []*Page `json:"pages"`
}

type Page struct {
	Title  string   `json:"title"`
	Source string   `json:"source"`
	Link   string   `json:"link"`
	Files  []string `json:"files"`
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
			renderPage(m, g, p)
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

func renderPage(mx *http.ServeMux, guide *Guide, page *Page) {
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
	mx.HandleFunc(page.Link, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	}))
	for _, name := range page.Files {
		file := filepath.Join(filepath.Dir(src), name)
		mx.HandleFunc(path.Join(path.Dir(page.Link), name), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, file)
		})
		mx.HandleFunc(path.Join(page.Link, name), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, file)
		})
	}
}
