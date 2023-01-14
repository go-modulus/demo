package event

type Event interface {
}

type Collector interface {
	PopEvents() []Event
}
