package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/vinceanalytics/sitetools/internal/copy"
	"github.com/vinceanalytics/sitetools/internal/svr"
)

var (
	eject = flag.String("eject", "", "")
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "blog":
		title := flag.Arg(1)
		if title != "" {
			svr.NewBlog(title)
		}
		return
	}
	h := svr.Hand()
	server := &http.Server{
		Addr:    ":9090",
		Handler: h,
	}
	if path := *eject; path != "" {
		go func() {
			server.ListenAndServe()
		}()
		time.Sleep(10 * time.Millisecond)
		err := copy.Copy(path)
		if err != nil {
			log.Fatalf("copy files %v", err)
		}
		for _, link := range svr.Links() {
			w, err := http.Get("http://localhost:9090/" + link)
			if err != nil {
				log.Fatal(err, link)
			}
			if w.StatusCode != http.StatusOK {
				log.Fatal(w.Status, link)
			}
			o := filepath.Join(path, link)
			if filepath.Ext(o) == "" {
				// if no extension it means it is a html file
				o = filepath.Join(o, "index.html")
			}
			os.MkdirAll(filepath.Dir(o), 0755)
			f, err := os.Create(o)
			if err != nil {
				log.Fatal(o, err)
			}
			f.ReadFrom(w.Body)
			f.Close()
			w.Body.Close()
			fmt.Println(link, "=>", o)
		}
		copyCharts(*svr.Root, path)
		return
	}
	server.ListenAndServe()
}

func copyCharts(in, out string) {
	src := filepath.Join(in, "charts")
	filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {

		rel, _ := filepath.Rel(in, path)
		dest := filepath.Join(out, rel)
		if info.IsDir() {
			os.MkdirAll(dest, 0755)
			return nil
		}
		x, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer x.Close()

		y, err := os.Open(path)
		if err != nil {
			return err
		}
		defer y.Close()
		x.ReadFrom(y)
		fmt.Println(path, "=>", dest)
		return nil
	})

}
