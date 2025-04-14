package root

import (
	"embed"
)

//go:embed *.tmpl
var Files embed.FS
