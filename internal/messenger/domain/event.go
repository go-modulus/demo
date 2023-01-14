package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type Event interface {
	Id() uuid.UUID
	OccurredOn() time.Time
}

type BaseEvent struct {
	id         uuid.UUID
	occurredOn time.Time
}

func (e *BaseEvent) Id() uuid.UUID {
	return e.id
}

func (e *BaseEvent) OccurredOn() time.Time {
	return e.occurredOn
}

type EventCollector struct {
	events []Event
}

func (c *EventCollector) recordThat(event Event) {
	c.events = append(c.events, event)
}

func (c *EventCollector) PopEvents() []Event {
	events := c.events
	c.events = nil
	return events
}
