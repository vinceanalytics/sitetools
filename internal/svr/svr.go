package svr

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/vinceanalytics/sitetools/data"
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

func Home(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFS(data.Templates, "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, map[string]any{})
}
