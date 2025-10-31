//go:build !dev

package ui

import (
	"embed"

	"github.com/pocketbase/pocketbase/apis"
)

// DistDirFS exposes the embedded Nuxt build output (without the "dist" prefix).
//
//go:embed all:dist
var dist embed.FS

var DistDirFS = apis.MustSubFS(dist, "dist")
