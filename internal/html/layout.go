package html

import (
	"boilerplate/internal/framework"
	"embed"
	"html/template"
	"net/http"
)

type LayoutBlock string

const (
	LayoutBlockContent     LayoutBlock = "content"
	LayoutBlockTitle       LayoutBlock = "title"
	LayoutBlockCurrentUser LayoutBlock = "currentUser"
	LayoutBlockErrors      LayoutBlock = "errors"
)

func (b LayoutBlock) String() string {
	return string(b)
}

//go:embed template
var LayoutFolder embed.FS

var indexLayout = template.Must(
	template.ParseFS(LayoutFolder, "template/index.gohtml"),
)

var ajaxLayout = template.Must(
	template.ParseFS(LayoutFolder, "template/ajax_content.gohtml"),
)

var errTemplate = template.Must(
	template.ParseFS(LayoutFolder, "template/errors.gohtml"),
)

type IndexPage interface {
	framework.Layout
}

type AjaxPage interface {
	framework.Layout
}

func NewIndexPage() IndexPage {
	headers := http.Header{}

	headers.Set("Content-Type", "text/html; charset=utf-8")
	return framework.NewPage(indexLayout).
		WithBlocks(errTemplate.Templates()).
		WithDefaultHeaders(headers)
}

func NewAjaxPage() AjaxPage {
	headers := http.Header{}
	headers.Set("Location", "/ajax/users")
	headers.Set("Content-Type", "text/vnd.turbo-stream.html")
	return framework.NewPage(ajaxLayout).
		WithBlocks(errTemplate.Templates()).
		WithDefaultHeaders(headers)
}
