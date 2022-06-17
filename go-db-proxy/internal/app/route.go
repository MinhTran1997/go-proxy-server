package app

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET  = "GET"
	POST = "POST"
)

func Route(r *mux.Router, ctx context.Context, conf Root) (*ApplicationContext, error) {
	app, err := NewApp(ctx, conf)
	if err != nil {
		return app, err
	}
	Handle(r, "/health", app.HealthHandler.Check, GET)
	Handle(r, "/begin", app.Handler.BeginTransaction, GET, POST)
	Handle(r, "/end", app.Handler.EndTransaction, GET, POST)
	Handle(r, "/query", app.Handler.Query, POST)
	Handle(r, "/exec", app.Handler.Exec, POST)
	Handle(r, "/exec-batch", app.Handler.ExecBatch, POST)
	return app, nil
}

func Handle(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), methods ...string) *mux.Route {
	return r.HandleFunc(path, f).Methods(methods...)
}
