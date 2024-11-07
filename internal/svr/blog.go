package svr

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gosimple/slug"
)

var blog []*Post

type Post struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Link   string `json:"link"`
	Author struct {
		Name   string `json:"name"`
		Social string `json:"social"`
	} `json:"author"`
	Date    Date          `json:"date"`
	Content template.HTML `json:"-"`
	Exerpt  template.HTML `json:"-"`
	Summary string        `json:"-"`
}

type Date struct {
	time.Time
}

func (ts *Date) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	ts.Time, err = time.Parse(time.DateOnly, s)
	return err
}

func (ts *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.Format(time.DateOnly))
}

func (p *Post) FormatDate() string {
	return p.Date.Format("January 02, 2006")
}

func NewBlog(title string) {
	blogPath := filepath.Join(*Root, "blog.json")
	data, err := os.ReadFile(blogPath)
	if err != nil {
		log.Fatal("reading blog data", err)
	}
	err = json.Unmarshal(data, &blog)
	if err != nil {
		log.Fatal("decoding blog data", err)
	}
	base := slug.Make(title)
	source := filepath.Join(filepath.Dir(blogPath), "blog", base+".md")
	os.MkdirAll(filepath.Dir(source), 0755)
	err = os.WriteFile(source, []byte{}, 0600)
	if err != nil {
		log.Fatal("creating blog file", err)
	}
	src, _ := filepath.Rel(*Root, source)
	blog = append([]*Post{
		{
			Title:  title,
			Source: src,
			Link:   "/blog/" + base,
			Author: struct {
				Name   string "json:\"name\""
				Social string "json:\"social\""
			}{
				Name:   "Geofrey Ernest",
				Social: "https://github.com/gernest",
			},
			Date: Date{Time: time.Now()},
		},
	}, blog...)
	dt, _ := json.MarshalIndent(blog, "", "  ")
	err = os.WriteFile(blogPath, dt, 0600)
	if err != nil {
		log.Fatal("updating blog metadata", err)
	}
}

func registerBlog(m *http.ServeMux) {
	blogPath := filepath.Join(*Root, "blog.json")
	data, err := os.ReadFile(blogPath)
	if err != nil {
		log.Fatal("reading blog data", err)
	}
	err = json.Unmarshal(data, &blog)
	if err != nil {
		log.Fatal("decoding blog data", err)
	}
	for _, p := range blog {
		src := filepath.Join(*Root, p.Source)
		md, err := os.ReadFile(src)
		if err != nil {
			log.Fatal(src, err)
		}
		p.Content = template.HTML(markdown.ToHTML(md, nil, nil))
		ex, _, _ := bytes.Cut(md, []byte("\n\n"))
		p.Summary = string(ex)
		p.Exerpt = template.HTML(markdown.ToHTML(ex, nil, nil))
		m.HandleFunc(p.Link, renderPost(p))
	}
	m.HandleFunc("/blog", renderBlog())
}

func renderPost(page *Post) http.HandlerFunc {
	var o bytes.Buffer
	err := global.ExecuteTemplate(&o, "post", baseContext(func(m map[string]any) {
		m["post"] = page
		m["title"] = page.Title + " | Vince Blog "
		m["meta"] = []Meta{
			{"description", page.Summary},
			{"og:site_name", "Vince"},
			{"og:title", page.Title},
			{"og:description", page.Summary},
			{"og:url", "https://vinceanalytics.com" + page.Link},
			{"og:image", "https://vinceanalytics.com/images/logo.png"},
			{"og:locale", "en_US"},
			{"og:type", "article"},
			{"twitter:site", "@gernesti"},
			{"twitter:card", "summary_large_image"},
		}
	}))
	if err != nil {
		log.Fatal("rendering post template", page.Link, err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
}

func renderBlog() http.HandlerFunc {
	var o bytes.Buffer
	err := global.ExecuteTemplate(&o, "blog", baseContext(func(m map[string]any) {
		m["blog"] = blog
		m["title"] = "Blog | Vince"
		m["meta"] = []Meta{
			{"description", "Blog posts about vince the cloud native web analytics server"},
			{"og:site_name", "Vince"},
			{"og:title", "Blog | Vince"},
			{"og:url", "https://vinceanalytics.com/blog"},
			{"og:image", "https://vinceanalytics.com/images/logo.png"},
			{"og:locale", "en_US"},
			{"og:type", "website"},
			{"twitter:site", "@gernesti"},
			{"twitter:card", "summary_large_image"},
		}
	}))
	if err != nil {
		log.Fatal("rendering blog template", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
}
