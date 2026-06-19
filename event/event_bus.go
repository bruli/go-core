package event

import (
	"context"
	"errors"
	"fmt"
)

type UnknownEventToDispatchError struct {
	event string
}

func (u UnknownEventToDispatchError) Error() string {
	return fmt.Sprintf("event %q is not declared to dispatch", u.event)
}

type Listener interface {
	Listen(ctx context.Context, ev Event) error
}

type Bus map[string][]Listener

func NewBus() Bus {
	return make(map[string][]Listener)
}

func (e Bus) Subscribe(eventName string, listeners ...Listener) {
	e[eventName] = listeners
}

func (e Bus) Dispatch(ctx context.Context, ev Event) error {
	errs := make([]error, 0)
	list, ok := e[ev.EventName()]
	if !ok {
		return UnknownEventToDispatchError{event: ev.EventName()}
	}
	for _, l := range list {
		if err := l.Listen(ctx, ev); err != nil {
			errs = append(errs, err)
		}
	}
	err := errors.Join(errs...)
	if err != nil {
		return NewDispatchError(err.Error(), ev.EventName())
	}
	return nil
}

type DispatchError struct {
	msg, eventName string
}

func (d DispatchError) Error() string {
	return fmt.Sprintf("failed to dispatch event %s error: %s", d.eventName, d.msg)
}

func NewDispatchError(msg string, eventName string) DispatchError {
	return DispatchError{msg: msg, eventName: eventName}
}
