package pkg

import "embed"

//go:embed services/*
var EmbeddedFiles embed.FS
