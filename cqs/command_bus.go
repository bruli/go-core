package cqs

import (
	"context"
	"fmt"

	"github.com/bruli-lab/go-core/event"
)

type CommandBus struct {
	m map[string]CommandHandler
}

func (c CommandBus) Handle(ctx context.Context, cmd Command) ([]event.Event, error) {
	hand, ok := c.m[cmd.Name()]
	if !ok {
		return nil, UnSubscribedCommandError{name: cmd.Name()}
	}
	return hand.Handle(ctx, cmd)
}

func NewCommandBus() CommandBus {
	m := make(map[string]CommandHandler)
	return CommandBus{m: m}
}

func (c CommandBus) Subscribe(name string, command CommandHandler) {
	c.m[name] = command
}

type UnSubscribedCommandError struct {
	name string
}

func (u UnSubscribedCommandError) Error() string {
	return fmt.Sprintf("command %q not subscribed", u.name)
}
