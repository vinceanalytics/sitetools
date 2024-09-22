package svr

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
