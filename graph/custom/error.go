package custom

import (
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
)

type ErrorPlugin struct {
}

var (
	_ plugin.Plugin              = &ErrorPlugin{}
	_ plugin.EarlySourceInjector = &ErrorPlugin{}
)

func (e ErrorPlugin) Name() string {
	//TODO implement me
	panic("implement me")
}

func (e ErrorPlugin) InjectSourceEarly() *ast.Source {
	//TODO implement me
	panic("implement me")
}
