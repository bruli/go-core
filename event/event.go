package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	EventID() uuid.UUID
	EventName() string
	EventAt() time.Time
	AggregateRootID() string
}
type BasicEvent struct {
	IDAttr              uuid.UUID `json:"id"`
	NameAttr            string    `json:"name"`
	AtAttr              time.Time `json:"at"`
	AggregateRootIDAttr string    `json:"aggregate_root_id"`
}

func (b BasicEvent) EventID() uuid.UUID {
	return b.IDAttr
}

func (b BasicEvent) EventName() string {
	return b.NameAttr
}

func (b BasicEvent) EventAt() time.Time {
	return b.AtAttr
}

func (b BasicEvent) AggregateRootID() string {
	return b.AggregateRootIDAttr
}

func NewBasicEvent(name string, id uuid.UUID, aggRootID string) BasicEvent {
	return BasicEvent{
		IDAttr:              id,
		NameAttr:            name,
		AtAttr:              time.Now(),
		AggregateRootIDAttr: aggRootID,
	}
}

type BasicAggregateRoot struct {
	createdAt time.Time
	events    []Event
}

func (b *BasicAggregateRoot) CreatedAt() time.Time {
	return b.createdAt
}

func NewBasicAggregateRoot() BasicAggregateRoot {
	return BasicAggregateRoot{
		createdAt: time.Now(),
		events:    nil,
	}
}

func (b *BasicAggregateRoot) Record(evs ...Event) {
	b.events = append(b.events, evs...)
}

func (b *BasicAggregateRoot) Events() []Event {
	events := b.events
	b.ClearEvents()
	return events
}

func (b *BasicAggregateRoot) ClearEvents() {
	b.events = nil
}

func (b *BasicAggregateRoot) Hydrate(createdAt time.Time, events []Event) {
	b.createdAt = createdAt
	b.events = events
}
