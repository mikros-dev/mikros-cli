package service_repository

import (
	"embed"
)

//go:embed assets/root/*.tmpl
var rootFiles embed.FS

//go:embed assets/scripts/*.tmpl
var scriptFiles embed.FS
