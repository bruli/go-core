package cqs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli-lab/go-core/cqs"
	"github.com/stretchr/testify/require"
)

func TestQueryBusHandle(t *testing.T) {
	errTest := errors.New("")
	tests := []struct {
		name, queryName string
		handler         cqs.QueryHandler
		query           cqs.Query
		expectedResult  any
		expectedErr     error
	}{
		{
			name:        "with a not subscribed query, then it returns an unsubscribed query error",
			queryName:   "query",
			handler:     queryHandler{},
			query:       query{name: "unknown"},
			expectedErr: cqs.UnSubscribedQueryError{},
		},
		{
			name:      "with a subscribed query, then it execute handle method",
			queryName: "query",
			handler:   queryHandler{},
			query:     query{name: "query"},
		},
		{
			name:      "with a subscribed query, then it execute handle method and return same query error",
			queryName: "other query",
			handler: queryHandler{
				err: errTest,
			},
			query:       query{name: "other query"},
			expectedErr: errTest,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a QueryBus,
		when Handle method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			bus := cqs.NewQueryBus()
			bus.Subscribe(tt.queryName, tt.handler)
			result, err := bus.Handle(context.Background(), tt.query)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				require.Nil(t, result)
				return
			}
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

type query struct {
	name string
}

func (c query) Name() string {
	return c.name
}

type queryHandler struct {
	result any
	err    error
}

func (q queryHandler) Handle(_ context.Context, _ cqs.Query) (any, error) {
	return q.result, q.err
}
