package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/vinceanalytics/sitetools/internal/build"
	assets "github.com/vinceanalytics/sitetools/internal/copy"
	"github.com/vinceanalytics/sitetools/internal/render"
)

func main() {
	s := flag.Bool("s", false, "")
	flag.Parse()
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

	_ = idx
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
