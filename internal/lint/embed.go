package lint

import (
	"embed"
)

//go:embed config/revive.toml.tmpl
var reviveConfig embed.FS
