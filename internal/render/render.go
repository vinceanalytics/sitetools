package render

import (
	"html/template"

	"github.com/vinceanalytics/sitetools/data"
)

type Render map[string]*template.Template

func Create() (pages, index Render) {
	pages = make(Render)
	index = make(Render)

	base := template.Must(template.ParseFS(
		data.Templates, "templates/include/include.html",
	))

	clone := func(name string) *template.Template {
		b := template.Must(base.Lookup("default").Clone())
		return template.Must(
			b.ParseFS(data.Templates, name),
		)
	}
	pages["article"] = clone("templates/layout/article.html")
	pages["page"] = clone("templates/layout/page.html")
	pages["post"] = clone("templates/layout/post.html")

	index["article"] = clone("templates/index/article.html")
	index["page"] = clone("templates/index/page.html")
	index["post"] = clone("templates/index/post.html")
	index["start"] = clone("templates/index/start.html")
	return
}
