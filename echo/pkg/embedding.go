package pkg

import "embed"

//go:embed services/base.go
//go:embed services/authentication/definition.go
//go:embed services/carpay/definition.go
var EmbeddedFiles embed.FS
