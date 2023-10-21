package html

import (
	"boilerplate/internal/auth/widget"
	"boilerplate/internal/framework"
	template2 "boilerplate/internal/html/template"
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

var errTemplate = template.Must(
	template.ParseFS(LayoutFolder, "template/errors.gohtml"),
)

type IndexPage interface {
	framework.Layout
}

type AjaxPage interface {
	framework.Layout
}

func NewIndexPage(
	currentUserWidget widget.CurrentUserWidget,
	config *ModuleConfig,
) IndexPage {
	headers := http.Header{}

	headers.Set("Content-Type", "text/html; charset=utf-8")
	return framework.NewPage(
		"index.gohtml",
		template2.GetTplFs(config.EmbeddedTemplates),
		config.UseCache,
	).
		WithBlocks(errTemplate.Templates()).
		WithWidget(currentUserWidget, []string{LayoutBlockCurrentUser.String()}).
		WithDefaultHeaders(headers)
}

func NewAjaxPage(
	config *ModuleConfig,
) AjaxPage {
	headers := http.Header{}
	headers.Set("Location", "/ajax/users")
	headers.Set("Content-Type", "text/vnd.turbo-stream.html")
	return framework.NewPage(
		"ajax_content.gohtml",
		template2.GetTplFs(config.EmbeddedTemplates),
		config.UseCache,
	).
		WithBlocks(errTemplate.Templates()).
		WithDefaultHeaders(headers)
}
