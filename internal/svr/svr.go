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

type Feature struct {
	Link  string
	Title string
	Body  string
}

func Home(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFS(data.Templates, "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/html")
	err = tpl.Execute(w, map[string]any{
		"features": []Feature{
			{
				Link:  "/guide/deployment/local",
				Title: "Self hosted",
				Body:  "Designed from grounds up for painless self hosting.",
			},
			{
				Link:  "/guide/dashboard/filters",
				Title: "Powerful filters",
				Body:  "Easily filter youd data to extract valuable insights",
			},
			{
				Link:  "/guide/dashboard/time-period",
				Title: "Time Period Comparison",
				Body:  "Compare data across different time periods for trend analysis.",
			},
			{
				Link:  "/guide/dashboard/session",
				Title: "Session Analysis",
				Body:  "Learn more about individual user journeys with in-depth session summaries.",
			},
			{
				Link:  "/guide/dashboard/custom-event",
				Title: "Custom Event Tracking",
				Body:  "Track and analyze custom events tailored to your website's needs.",
			},
			{
				Link:  "/guide/dashboard/404",
				Title: "404 Page Tracking",
				Body:  "Identify and address broken links with 404 page tracking.",
			},
		},
	})
	if err != nil {
		log.Println(err)
	}
}
