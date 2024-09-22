package svr

import (
	"bytes"
	"html/template"
	"log"

	"github.com/vinceanalytics/sitetools/data"
)

func footer() template.HTML {
	tpl, err := template.ParseFS(data.Templates, "templates/footer.html")
	if err != nil {
		log.Fatal("parsing footer template", err)
	}
	var o bytes.Buffer
	err = tpl.Execute(&o, map[string]any{
		"guides": guideIndex(),
	})
	if err != nil {
		log.Fatal("rendering footer template", err)
	}
	return template.HTML(o.String())
}
