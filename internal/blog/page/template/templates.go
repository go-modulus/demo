package template

import (
	"boilerplate/internal/framework"
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
	return os.DirFS("internal/blog/page/template")
}

func GetNewPost(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("new_post.gohtml", tplFs)
}

func GetPosts(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("posts.gohtml", tplFs)
}

func GetPost(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("post.gohtml", tplFs)
}
