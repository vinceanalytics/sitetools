package main

import (
	"flag"
	"net/http"

	"github.com/vinceanalytics/sitetools/internal/svr"
)

func main() {
	flag.Parse()
	h := svr.Hand()
	server := &http.Server{
		Addr:    ":9090",
		Handler: h,
	}
	server.ListenAndServe()
}
