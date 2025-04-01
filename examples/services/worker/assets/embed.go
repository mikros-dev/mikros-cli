<<<<<<<< HEAD:internal/assets/golang/embed.go
package golang
========
package assets
>>>>>>>> main:examples/services/worker/assets/embed.go

import (
	"embed"
)

//go:embed *.tmpl
var Files embed.FS
