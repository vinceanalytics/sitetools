package svr

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vinceanalytics/sitetools/data"
)

var (
	root = flag.String("root", "", "")
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
	return h
}

type Feature struct {
	Link  string `json:"link"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFS(data.Templates, "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/html")
	var f []Feature
	featureFile := filepath.Join(*root, "features.json")
	file, err := os.ReadFile(featureFile)
	if err != nil {
		log.Fatalf("failed reading features file %s %v", featureFile, err)
	}
	err = json.Unmarshal(file, &f)
	if err != nil {
		log.Fatalf("failed decoding features file %s %v", featureFile, err)
	}
	err = tpl.Execute(w, map[string]any{
		"features": f,
	})
	if err != nil {
		log.Println(err)
	}
}
