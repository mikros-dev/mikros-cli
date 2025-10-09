package protobuf_module

import (
	"embed"
)

//go:embed assets/*.tmpl
var templateFiles embed.FS
