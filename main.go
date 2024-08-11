package main

import (
	"flag"
	"log/slog"
	"os"

	assets "github.com/vinceanalytics/sitetools/internal/copy"
)

func main() {
	flag.Parse()
	dst := flag.Arg(0)
	err := assets.Copy(dst)
	if err != nil {
		slog.Error("copying assets", "err", err)
		os.Exit(1)
	}
}
