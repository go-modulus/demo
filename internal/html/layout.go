package html

import (
	"boilerplate/internal/auth/widget"
	"boilerplate/internal/framework"
	"embed"
	"html/template"
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

func NewIndexPage(
	currentUserWidget widget.CurrentUserWidget,
) (IndexPage, error) {
	return framework.NewPage(indexLayout).
		WithBlocks(errTemplate.Templates()).
		WithWidget(currentUserWidget), nil
}

func NewAjaxPage() AjaxPage {
	return framework.NewPage(ajaxLayout)
}
