package override

import (
	"boilerplate/internal/user/httpaction"
	"go.uber.org/fx"
)

type TestOverride struct {
}

func NewTestOverride() *TestOverride {
	return &TestOverride{}
}

func (t TestOverride) Name() string {
	return "OVER"
}
func Overrides() fx.Option {
	return fx.Decorate(
		func() httpaction.TestOverride {
			return &TestOverride{}
		},
	)
}
