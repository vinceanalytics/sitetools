package main

import (
	"net/http"

	"github.com/vinceanalytics/sitetools/internal/svr"
)

func main() {
	h := svr.Hand()
	server := &http.Server{
		Addr:    ":9090",
		Handler: h,
	}
	server.ListenAndServe()
}
