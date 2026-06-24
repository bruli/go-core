package cqs

import (
	"context"

	"github.com/bruli/go-core/event"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewCommandHandlerEventListenerMiddleware(list event.Listener, tracer trace.Tracer) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, c Command) ([]event.Event, error) {
			ctx, span := tracer.Start(ctx, "CommandHandlerEventListenerMiddleware")
			defer span.End()
			events, err := h.Handle(ctx, c)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return nil, err
			}
			for _, ev := range events {
				if err := list.Listen(ctx, ev); err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					return nil, err
				}
			}
			span.SetStatus(codes.Ok, "event listener middleware completed")
			return events, nil
		})
	}
}
