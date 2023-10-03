package framework

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type Layout interface {
	WithWidget(widget Widget) *Page
	Handler(
		successCode int,
		successHeaders http.Header,
		errorHeader http.Header,
	) http.HandlerFunc
}

type Widget interface {
	Blocks() []*template.Template
	DataSource() PageDataSource
	DataHolderNames() []string
}

type BaseWidget struct {
	blocks          []*template.Template
	dataSource      PageDataSource
	dataHolderNames []string
}

func NewWidget(template *template.Template, dataSource PageDataSource, dataHolderNames []string) *BaseWidget {
	return &BaseWidget{blocks: template.Templates(), dataSource: dataSource, dataHolderNames: dataHolderNames}
}

func (w *BaseWidget) Blocks() []*template.Template {
	return w.blocks
}

func (w *BaseWidget) DataSource() PageDataSource {
	return w.dataSource
}

func (w *BaseWidget) DataHolderNames() []string {
	return w.dataHolderNames
}

type Page struct {
	layout                   *template.Template
	errorHandler             PageErrorHandler
	dataSources              map[string]PageDataSource
	dataHoldersToHandlersMap map[string]string
	errorVarName             string
	defaultHeaders           http.Header
}

func NewPage(layout *template.Template) *Page {
	return &Page{
		layout:                   layout,
		dataSources:              make(map[string]PageDataSource),
		dataHoldersToHandlersMap: make(map[string]string),
		errorVarName:             "errors",
		errorHandler: func(w http.ResponseWriter, req *http.Request, errors []error) []error {
			return errors
		},
	}
}

func (p *Page) clone() *Page {
	ds := make(map[string]PageDataSource)
	for k, v := range p.dataSources {
		ds[k] = v
	}
	dh := make(map[string]string)
	for k, v := range p.dataHoldersToHandlersMap {
		dh[k] = v
	}
	return &Page{
		layout:                   template.Must(p.layout.Clone()),
		dataSources:              ds,
		dataHoldersToHandlersMap: dh,
		errorVarName:             p.errorVarName,
		errorHandler:             p.errorHandler,
		defaultHeaders:           p.defaultHeaders,
	}
}

// addBlocks adds blocks to the layout template
// If the block already exists in the layout, it will be appended to the existing block
// If the block does not exist in the layout, it will be added to the layout
// Use it only for blocks that reused in every page
// Use WithBlocks instead for a block that has uniq content for each page
func (p *Page) addBlocks(blocks []*template.Template) error {
	layoutBlocksMap := make(map[string]*template.Template)

	for _, el := range p.layout.Templates() {
		layoutBlocksMap[el.Name()] = el
	}

	for _, el := range blocks {
		if block, ok := layoutBlocksMap[el.Name()]; ok {
			elComposite, err := block.Parse(block.Tree.Root.String() + el.Tree.Root.String())
			if err != nil {
				return err
			}
			el = elComposite

		}

		template.Must(p.layout.AddParseTree(el.Name(), el.Tree))
	}
	return nil
}

// WithBlocks adds blocks to the layout template
// Returns a new Page with the new layout
func (p *Page) WithBlocks(blocks []*template.Template) *Page {
	newPage := p.clone()
	err := newPage.addBlocks(blocks)
	if err != nil {
		panic(err)
	}
	return newPage
}

// WithDefaultHeaders adds default headers to the responses for this page
// Returns a new Page with the new layout
func (p *Page) WithDefaultHeaders(defaultHeaders http.Header) *Page {
	newPage := p.clone()
	newPage.defaultHeaders = defaultHeaders
	return newPage
}

// setDataSource sets the data source for the page
// The data source will be used to populate the variables in layout to be sent to the blocks
// For example:
// <div>{{ block "content" .dataHolder }}{{ end }}</div>
// dataHolder is the data holder name. It will be populated with the data from the data source
// Use WithDataSource instead to return a new Page with the new data source
func (p *Page) setDataSource(
	datasource PageDataSource,
	dataHolderNames []string,
) {
	dsName := ""
	for _, el := range dataHolderNames {
		dsName += el
	}
	p.dataSources[dsName] = datasource
	for _, el := range dataHolderNames {
		p.dataHoldersToHandlersMap[el] = dsName
	}
}

func (p *Page) WithDataSource(
	datasource PageDataSource,
	dataHolderNames []string,
) *Page {
	newPage := p.clone()
	newPage.setDataSource(datasource, dataHolderNames)
	return newPage
}

// WithWidget adds blocks and data sources to the page to visualise a prt of content
// Returns a new Page with the new layout
func (p *Page) WithWidget(
	widget Widget,
) *Page {
	newPage := p.clone()
	err := newPage.addBlocks(widget.Blocks())
	if err != nil {
		panic(err)
	}
	newPage.setDataSource(widget.DataSource(), widget.DataHolderNames())

	return newPage
}

func (p *Page) SetErrorVarName(name string) {
	p.errorVarName = name
}

func (p *Page) SetErrorHandler(errorHandler PageErrorHandler) {
	p.errorHandler = errorHandler
}

func (p *Page) Handler(
	successCode int,
	successHeaders http.Header,
	errorHeader http.Header,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		data, _ := io.ReadAll(req.Body)
		req.Body = RequestBody{bytes.NewReader(data)}
		ctx = SetHttpRequest(ctx, req)
		ctx = SetHttpResponseWriter(ctx, w)
		req = req.WithContext(ctx)

		tpl := template.Must(p.layout.Clone())

		tplData := make(map[string]any)
		errors := make([]error, 0)
		handlersResultCache := make(map[string]any)
		for varName, handlerName := range p.dataHoldersToHandlersMap {
			var res any
			if cachedRes, ok := handlersResultCache[handlerName]; ok {
				res = cachedRes
			} else {
				var widgetErr error
				if handler, ok := p.dataSources[handlerName]; ok {
					res, widgetErr = handler(w, req)
					handlersResultCache[handlerName] = res
					if widgetErr != nil {
						errors = append(errors, widgetErr)
					}
				}
			}

			tplData[varName] = res
		}

		tplData[p.errorVarName] = []error{}
		for key, headers := range p.defaultHeaders {
			for _, header := range headers {
				w.Header().Set(key, header)
			}
		}

		if len(errors) > 0 {
			if p.errorHandler != nil {
				errors = p.errorHandler(w, req, errors)
			}
			tplData[p.errorVarName] = errors

			if errorHeader != nil {
				for key, headers := range errorHeader {
					for _, header := range headers {
						w.Header().Set(key, header)
					}
				}
			}
		} else {
			if successHeaders != nil {
				for key, headers := range successHeaders {
					for _, header := range headers {
						w.Header().Set(key, header)
					}
				}
			}
			w.Header().Set(
				"Status Code",
				fmt.Sprintf("%d %s", successCode, http.StatusText(successCode)),
			)
			w.WriteHeader(successCode)
		}

		err := tpl.Execute(w, tplData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(successCode)
		}
	}
}
