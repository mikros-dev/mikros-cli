package service

import (
	"embed"
)

//go:embed assets/*.tmpl
var templateFiles embed.FS
