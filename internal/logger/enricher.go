package logger

import "context"

type Enricher interface {
	Enrich(ctx context.Context) map[string]string
}

type EnricherFunc func(context.Context) map[string]string

func (e EnricherFunc) Enrich(ctx context.Context) map[string]string {
	return e(ctx)
}

type RootEnricher struct {
	Enrichers []EnricherFunc
}

func NewRootEnricher() *RootEnricher {
	return &RootEnricher{}
}

func (e *RootEnricher) AttachEnricher(ef EnricherFunc) {
	e.Enrichers = append(e.Enrichers, ef)
}

func (e *RootEnricher) Enrich(ctx context.Context) map[string]string {
	fields := make(map[string]string)

	for _, e := range e.Enrichers {
		for k, v := range e(ctx) {
			fields[k] = v
		}
	}

	return fields
}
