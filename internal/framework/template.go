package framework

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"sync"
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
	DataSource() *PageDataSource
}

type BaseWidget struct {
	templatePath string
	fs           fs.FS
	dataSource   *PageDataSource
}

func NewWidget(
	templatePath string,
	fs fs.FS,
	dataSource *PageDataSource,
) *BaseWidget {
	return &BaseWidget{
		templatePath: templatePath,
		fs:           fs,
		dataSource:   dataSource,
	}
}

func (w *BaseWidget) Blocks() []*template.Template {
	return template.Must(
		template.ParseFS(w.fs, w.templatePath),
	).Templates()
}

func (w *BaseWidget) DataSource() *PageDataSource {
	return w.dataSource
}

type Page struct {
	layout         *template.Template
	templatePath   string
	fs             fs.FS
	useCache       bool
	errorHandler   PageErrorHandler
	dataSources    []*PageDataSource
	errorVarName   string
	defaultHeaders http.Header
	widgets        []Widget
	mu             sync.Mutex
}

func NewPage(
	templatePath string,
	fs fs.FS,
	useCache bool,
) *Page {
	_, err := fs.Open(templatePath)
	if err != nil {
		panic("template file not found in the fs " + templatePath)
	}
	return &Page{
		templatePath: templatePath,
		fs:           fs,
		dataSources:  make([]*PageDataSource, 0),
		useCache:     useCache,
		errorVarName: "errors",
		errorHandler: func(w http.ResponseWriter, req *http.Request, errors []error) []error {
			return errors
		},
		widgets: make([]Widget, 0),
	}
}

func (p *Page) clone() *Page {
	ds := make([]*PageDataSource, len(p.dataSources))
	for k, v := range p.dataSources {
		ds[k] = v
	}

	widgetsCopy := make([]Widget, len(p.widgets))
	copy(widgetsCopy, p.widgets)

	var layout *template.Template
	if p.layout != nil {
		layout = template.Must(p.layout.Clone())
	}
	return &Page{
		layout:         layout,
		templatePath:   p.templatePath,
		fs:             p.fs,
		useCache:       p.useCache,
		errorHandler:   p.errorHandler,
		dataSources:    ds,
		errorVarName:   p.errorVarName,
		defaultHeaders: p.defaultHeaders,
		widgets:        widgetsCopy,
		mu:             sync.Mutex{},
	}
}

func (p *Page) addWidget(widget Widget) error {
	p.widgets = append(p.widgets, widget)
	return nil
}

func (p *Page) parseBlocks(blocks []*template.Template) error {
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

		layoutBlocksMap[el.Name()] = el
	}

	for _, el := range layoutBlocksMap {
		template.Must(p.layout.AddParseTree(el.Name(), el.Tree))
	}
	return nil
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
	datasource *PageDataSource,
) {
	p.dataSources = append(
		p.dataSources,
		datasource,
	)
}

func (p *Page) WithDataSource(
	datasource *PageDataSource,
) *Page {
	newPage := p.clone()
	newPage.setDataSource(datasource)
	return newPage
}

// WithWidget adds blocks and data sources to the page to visualise a prt of content
// Returns a new Page with the new layout
func (p *Page) WithWidget(
	widget Widget,
) *Page {
	newPage := p.clone()
	err := newPage.addWidget(widget)
	if err != nil {
		panic(err)
	}
	if widget.DataSource() != nil {
		newPage.setDataSource(widget.DataSource())
	}

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

		if !p.useCache || p.layout == nil {
			err := p.fillLayout()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}
		tpl := template.Must(p.layout.Clone())

		tplData := make(map[string]any)
		errors := make([]error, 0)
		wg := sync.WaitGroup{}
		wg.Add(len(p.dataSources))
		mu := sync.Mutex{}
		for _, source := range p.dataSources {
			go func(info *PageDataSource) {
				defer wg.Done()
				res, widgetErr := info.Handle(w, req)
				mu.Lock()
				defer mu.Unlock()
				if widgetErr != nil {
					errors = append(errors, widgetErr)
				}
				tplData[info.TplVarName()] = res
			}(source)
		}
		wg.Wait()

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
		}
	}
}

func (p *Page) fillLayout() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.layout = template.Must(template.ParseFS(p.fs, p.templatePath))
	blocks := make([]*template.Template, 0)

	for _, widget := range p.widgets {
		blocks = append(blocks, widget.Blocks()...)
	}
	err := p.parseBlocks(blocks)
	return err
}
