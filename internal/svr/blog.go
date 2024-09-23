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

func (p *Post) FormatDate() string {
	return p.Date.Format("January 02, 2006")
}

func registerBlog(m *http.ServeMux) {
	blogPath := filepath.Join(*root, "blog.json")
	data, err := os.ReadFile(blogPath)
	if err != nil {
		log.Fatal("reading blog data", err)
	}
	err = json.Unmarshal(data, &blog)
	if err != nil {
		log.Fatal("decoding blog data", err)
	}
	for _, p := range blog {
		src := filepath.Join(*root, p.Source)
		md, err := os.ReadFile(src)
		if err != nil {
			log.Fatal(src, err)
		}
		p.Content = template.HTML(markdown.ToHTML(md, nil, nil))
		ex, _, _ := bytes.Cut(md, []byte("\n\n"))
		p.Exerpt = template.HTML(markdown.ToHTML(ex, nil, nil))
		m.HandleFunc(p.Link, renderPost(p))
	}
	m.HandleFunc("/blog", renderBlog())
}

func renderPost(page *Post) http.HandlerFunc {
	var o bytes.Buffer
	err := global.ExecuteTemplate(&o, "post", baseContext(func(m map[string]any) {
		m["post"] = page
		m["title"] = page.Title
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
		m["title"] = "vince | Blog"
	}))
	if err != nil {
		log.Fatal("rendering blog template", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write(o.Bytes())
	})
}
