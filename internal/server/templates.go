package server

import _ "embed"

//go:embed templates/index.html
var indexRaw string

//go:embed templates/client.go.tmpl
var goRaw string
