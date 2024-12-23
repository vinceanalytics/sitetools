package svr

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/vinceanalytics/sitetools/data"
)

var (
	Root   = flag.String("root", ".", "")
	global = template.Must(template.ParseFS(data.Templates, "templates/global.html"))
)

func Hand() *http.ServeMux {
	h := http.NewServeMux()
	h.HandleFunc("/{$}", Home)
	static, err := fs.Sub(data.Assets, "assets")
	if err != nil {
		log.Fatal(err)
	}
	staticFS := http.FileServerFS(static)
	h.Handle("/css/", staticFS)
	h.Handle("/images/", staticFS)
	h.Handle("/js/", staticFS)
	registerGuides(h)
	registerBlog(h)
	return h
}

type Feature struct {
	Link  string `json:"link"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	var f []Feature
	featureFile := filepath.Join(*Root, "features.json")
	file, err := os.ReadFile(featureFile)
	if err != nil {
		log.Fatalf("failed reading features file %s %v", featureFile, err)
	}
	err = json.Unmarshal(file, &f)
	if err != nil {
		log.Fatalf("failed decoding features file %s %v", featureFile, err)
	}
	err = global.ExecuteTemplate(w, "home", baseContext(func(m map[string]any) {
		m["features"] = f
	}))
	if err != nil {
		log.Println("rendering home page", err)
	}
}

func baseContext(f ...func(map[string]any)) map[string]any {
	a := map[string]any{
		"guides": guides,
	}
	for i := range f {
		f[i](a)
	}
	return a
}

type Meta struct {
	Name    string
	Content string
}

func Links() (o []string) {
	o = []string{"/", "/blog", "/guides"}
	for _, g := range guides {
		for _, p := range g.Pages {
			o = append(o, p.Link)
			for _, f := range p.Files {
				o = append(o, path.Join(p.Link, f))
			}
		}
	}
	for _, p := range blog {
		o = append(o, p.Link)
	}
	return
}
