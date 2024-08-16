package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/go-github/v63/github"
	"github.com/vinceanalytics/sitetools/internal/build"
	assets "github.com/vinceanalytics/sitetools/internal/copy"
	"github.com/vinceanalytics/sitetools/internal/render"
)

func main() {
	s := flag.Bool("s", false, "")
	rel := flag.Bool("r", false, "")
	flag.Parse()
	if *rel {
		releases()
	}
	src := flag.Arg(0)
	dst := flag.Arg(1)
	err := assets.Copy(dst)
	if err != nil {
		slog.Error("copying assets", "err", err)
		os.Exit(1)
	}
	r, idx := render.Create()
	pgx := build.Build(src)

	pgx.Render(dst, r)

	pgx.RenderSpecial(dst, build.SpecialPages(), idx)
	if *s {
		serve(dst)
	}
}

func serve(dir string) {
	fs := http.FileServer(http.Dir(dir))
	slog.Info("serving site", "dir", dir, "addr", 8080)
	svr := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			x := &wrap{ResponseWriter: w}
			fs.ServeHTTP(x, r)
			slog.Info(r.URL.Path, "method", r.Method, "code", http.StatusText(x.code))
		}),
	}
	svr.ListenAndServe()
}

type wrap struct {
	code int
	http.ResponseWriter
}

func (w *wrap) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func releases() {
	client := github.NewClient(nil)
	all, _, err := client.Repositories.ListReleases(
		context.Background(),
		"vinceanalytics",
		"vince",
		nil,
	)
	if err != nil {
		log.Fatal("fetching releases", err)
	}
	data, _ := json.Marshal(all)
	os.WriteFile("releases.json", data, 0600)
}
