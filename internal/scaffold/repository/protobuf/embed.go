package protobuf

import (
	"embed"
)

//go:embed assets/proto/*.tmpl
var protoTemplateFiles embed.FS

//go:embed assets/root/*.tmpl
var rootTemplateFiles embed.FS

//go:embed assets/scripts/*.tmpl
var scriptsTemplateFiles embed.FS
