//go:build dev

package ui

import (
	"embed"
	"io/fs"

	"github.com/pocketbase/pocketbase/apis"
)

//go:embed dev/*
var dev embed.FS

var DistDirFS fs.FS = apis.MustSubFS(dev, "dev")
