package event_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bruli-lab/go-core/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBus_Dispatch(t *testing.T) {
	type args struct {
		subscribedEventName string
		dispatchedEventName string
		listeners           []event.Listener
	}
	tests := []struct {
		name           string
		args           args
		expectedErr    error
		expectedErrMsg string
	}{
		{
			name: "with an unsubscribed event, then it returns all messages in the error",
			args: args{
				subscribedEventName: "test",
				dispatchedEventName: "other",
			},
			expectedErr: event.UnknownEventToDispatchError{},
		},
		{
			name: "with all listeners returns an error, then it returns all messages in the error",
			args: args{
				listeners: []event.Listener{
					listener{msgError: new("error1")},
					listener{msgError: new("error2")},
				},
				subscribedEventName: "test",
				dispatchedEventName: "test",
			},
			expectedErr:    event.DispatchError{},
			expectedErrMsg: "failed to dispatch event test error: error1\nerror2",
		},
		{
			name: "with one listener returns an error, and the other that returns nil, then it returns the first error",
			args: args{
				listeners: []event.Listener{
					listener{msgError: new("error1")},
					listener{},
				},
				subscribedEventName: "test",
				dispatchedEventName: "test",
			},
			expectedErr:    event.DispatchError{},
			expectedErrMsg: "failed to dispatch event test error: error1",
		},
		{
			name: "with one listener returns an returns nil, then it returns nil",
			args: args{
				listeners: []event.Listener{
					listener{},
					listener{},
				},
				subscribedEventName: "test",
				dispatchedEventName: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a event bus with subscribed listeners to event
		when Dispatch method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			bus := event.NewBus()
			ev := eventTest{name: tt.args.dispatchedEventName}
			bus.Subscribe(tt.args.subscribedEventName, tt.args.listeners...)
			err := bus.Dispatch(t.Context(), ev)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			if errors.As(err, &event.DispatchError{}) {
				require.Equal(t, tt.expectedErrMsg, err.Error())
				return
			}

			if tt.expectedErr != nil {
				require.Error(t, err, "expected error but got nil")
			}
		})
	}
}

type eventTest struct {
	name string
}

func (e eventTest) EventID() uuid.UUID {
	return uuid.New()
}

func (e eventTest) EventName() string {
	return e.name
}

func (e eventTest) EventAt() time.Time {
	return time.Now()
}

func (e eventTest) AggregateRootID() string {
	return uuid.NewString()
}

type listener struct {
	msgError *string
}

func (l listener) Listen(_ context.Context, _ event.Event) error {
	if l.msgError != nil {
		return errors.New(*l.msgError)
	}
	return nil
}
