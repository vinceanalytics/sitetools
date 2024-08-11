package data

import "embed"

//go:embed assets
var Assets embed.FS

//go:embed templates
var Templates embed.FS

//go:embed CNAME
var Cname []byte
