package cqs

import (
	"context"
	"fmt"
)

type QueryBus struct {
	m map[string]QueryHandler
}

func (c QueryBus) Handle(ctx context.Context, query Query) (any, error) {
	hand, ok := c.m[query.Name()]
	if !ok {
		return nil, UnSubscribedQueryError{name: query.Name()}
	}
	return hand.Handle(ctx, query)
}

func NewQueryBus() QueryBus {
	m := make(map[string]QueryHandler)
	return QueryBus{m: m}
}

func (c QueryBus) Subscribe(name string, query QueryHandler) {
	c.m[name] = query
}

type UnSubscribedQueryError struct {
	name string
}

func (u UnSubscribedQueryError) Error() string {
	return fmt.Sprintf("query %q not subscribed", u.name)
}
