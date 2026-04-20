package assets

import "embed"

//go:embed web/templates/*.gohtml web/static/css/* web/static/js/*
var FS embed.FS
