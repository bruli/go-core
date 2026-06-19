package cqs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/go-core/cqs"
	"github.com/bruli/go-core/event"
	"github.com/stretchr/testify/require"
)

func TestCommandBusHandle(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name, cmdName  string
		handler        cqs.CommandHandler
		cmd            cqs.Command
		expectedEvents []event.Event
		expectedErr    error
	}{
		{
			name:        "with a not subscribed command, then it returns an unsubscribed command error",
			cmdName:     "command",
			handler:     commandHandler{},
			cmd:         command{name: "unknown"},
			expectedErr: cqs.UnSubscribedCommandError{},
		},
		{
			name:    "with a subscribed command, then it execute handle method",
			cmdName: "command",
			handler: commandHandler{},
			cmd:     command{name: "command"},
		},
		{
			name:    "with a subscribed command, then it execute handle method and return same command error",
			cmdName: "other command",
			handler: commandHandler{
				err: errTest,
			},
			cmd:         command{name: "other command"},
			expectedErr: errTest,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a CommandBus,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			bus := cqs.NewCommandBus()
			bus.Subscribe(tt.cmdName, tt.handler)
			events, err := bus.Handle(context.Background(), tt.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				require.Nil(t, events)
				return
			}
			require.Equal(t, tt.expectedEvents, events)
		})
	}
}

type command struct {
	name string
}

func (c command) Name() string {
	return c.name
}

type commandHandler struct {
	events []event.Event
	err    error
}

func (c commandHandler) Handle(_ context.Context, _ cqs.Command) ([]event.Event, error) {
	return c.events, c.err
}
