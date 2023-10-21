package template

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed *
var tplFolder embed.FS

func GetTplFs(embedded bool) fs.FS {
	if embedded {
		return tplFolder
	}
	return os.DirFS("internal/auth/widget/template")
}
