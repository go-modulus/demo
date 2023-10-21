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
	WithWidget(widget Widget, dataHolderNames []string) *Page
	Handler(
		successCode int,
		successHeaders http.Header,
		errorHeader http.Header,
	) http.HandlerFunc
}

type Widget interface {
	Blocks(useCache bool) []*template.Template
	DataSource() PageDataSource
}

type BaseWidget struct {
	templatePath string
	fs           fs.FS
	blocks       []*template.Template
	dataSource   PageDataSource
	mu           sync.Mutex
}

func NewWidget(
	templatePath string,
	fs fs.FS,
	dataSource PageDataSource,
) *BaseWidget {
	return &BaseWidget{
		templatePath: templatePath,
		fs:           fs,
		dataSource:   dataSource,
	}
}

func (w *BaseWidget) Blocks(useCache bool) []*template.Template {
	reParse := !useCache
	if reParse || w.blocks == nil {
		w.mu.Lock()
		defer w.mu.Unlock()
		w.blocks = template.Must(
			template.ParseFS(w.fs, w.templatePath),
		).Templates()
	}
	return w.blocks
}

func (w *BaseWidget) DataSource() PageDataSource {
	return w.dataSource
}

type datasourceInfo struct {
	handler     PageDataSource
	dataHolders []string
}

type Page struct {
	layout                   *template.Template
	templatePath             string
	fs                       fs.FS
	useCache                 bool
	errorHandler             PageErrorHandler
	dataSources              []datasourceInfo
	dataHoldersToHandlersMap map[string]string
	errorVarName             string
	defaultHeaders           http.Header
	templateBlocks           []*template.Template
	widgets                  []Widget
	mu                       sync.Mutex
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
		templatePath:             templatePath,
		fs:                       fs,
		dataSources:              make([]datasourceInfo, 0),
		useCache:                 useCache,
		dataHoldersToHandlersMap: make(map[string]string),
		errorVarName:             "errors",
		errorHandler: func(w http.ResponseWriter, req *http.Request, errors []error) []error {
			return errors
		},
		templateBlocks: make([]*template.Template, 0),
		widgets:        make([]Widget, 0),
	}
}

func (p *Page) clone() *Page {
	ds := make([]datasourceInfo, len(p.dataSources))
	for k, v := range p.dataSources {
		ds[k] = v
	}
	dh := make(map[string]string)
	for k, v := range p.dataHoldersToHandlersMap {
		dh[k] = v
	}
	tplBlocks := make([]*template.Template, len(p.templateBlocks))
	copy(tplBlocks, p.templateBlocks)

	widgetsCopy := make([]Widget, len(p.widgets))
	copy(widgetsCopy, p.widgets)

	var layout *template.Template
	if p.layout != nil {
		layout = template.Must(p.layout.Clone())
	}
	return &Page{
		layout:                   layout,
		templatePath:             p.templatePath,
		fs:                       p.fs,
		useCache:                 p.useCache,
		errorHandler:             p.errorHandler,
		dataSources:              ds,
		dataHoldersToHandlersMap: dh,
		errorVarName:             p.errorVarName,
		defaultHeaders:           p.defaultHeaders,
		templateBlocks:           tplBlocks,
		widgets:                  widgetsCopy,
		mu:                       sync.Mutex{},
	}
}

// addBlocks adds blocks to the layout template
// If the block already exists in the layout, it will be appended to the existing block
// If the block does not exist in the layout, it will be added to the layout
// Use it only for blocks that reused in every page
// Use WithBlocks instead for a block that has uniq content for each page
func (p *Page) addBlocks(blocks []*template.Template) error {
	p.templateBlocks = append(p.templateBlocks, blocks...)
	return nil // p.parseBlocks(blocks)
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
			layoutBlocksMap[el.Name()] = el
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
	p.dataSources = append(
		p.dataSources, datasourceInfo{
			handler:     datasource,
			dataHolders: dataHolderNames,
		},
	)
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
	dataHolderNames []string,
) *Page {
	newPage := p.clone()
	err := newPage.addWidget(widget)
	if err != nil {
		panic(err)
	}
	newPage.setDataSource(widget.DataSource(), dataHolderNames)

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
		//datasourcesInfo := make(map[string]*datasourceInfo, len(p.dataSources))
		//for varName, handlerName := range p.dataHoldersToHandlersMap {
		//	if dsInfo, ok := datasourcesInfo[handlerName]; ok {
		//		dsInfo.dataHolders = append(dsInfo.dataHolders, varName)
		//		continue
		//	}
		//	dsInfo := datasourceInfo{
		//		handler:     p.dataSources[handlerName],
		//		dataHolders: []string{varName},
		//	}
		//	datasourcesInfo[handlerName] = &dsInfo
		//}
		wg := sync.WaitGroup{}
		wg.Add(len(p.dataSources))
		mu := sync.Mutex{}
		for _, info := range p.dataSources {
			go func(info datasourceInfo) {
				defer wg.Done()
				res, widgetErr := info.handler(w, req)
				mu.Lock()
				defer mu.Unlock()
				if widgetErr != nil {
					errors = append(errors, widgetErr)
				}
				for _, varName := range info.dataHolders {
					tplData[varName] = res
				}
			}(info)
		}
		wg.Wait()
		//handlersResultCache := make(map[string]any)
		//for varName, handlerName := range p.dataHoldersToHandlersMap {
		//	var res any
		//	if cachedRes, ok := handlersResultCache[handlerName]; ok {
		//		res = cachedRes
		//	} else {
		//		var widgetErr error
		//		if handler, ok := p.dataSources[handlerName]; ok {
		//			res, widgetErr = handler(w, req)
		//			handlersResultCache[handlerName] = res
		//			if widgetErr != nil {
		//				errors = append(errors, widgetErr)
		//			}
		//		}
		//	}
		//
		//	tplData[varName] = res
		//}

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
	if len(p.templateBlocks) > 0 {
		blocks = append(blocks, p.templateBlocks...)
	}
	for _, widget := range p.widgets {
		blocks = append(blocks, widget.Blocks(p.useCache)...)
	}
	err := p.parseBlocks(blocks)
	return err
}
