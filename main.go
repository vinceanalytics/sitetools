package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/vinceanalytics/sitetools/internal/copy"
	"github.com/vinceanalytics/sitetools/internal/svr"
)

var (
	eject = flag.String("eject", "", "")
)

func main() {
	flag.Parse()
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
		return
	}
	server.ListenAndServe()
}
