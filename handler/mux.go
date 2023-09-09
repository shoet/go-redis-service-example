package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewMux() http.Handler {
	router := chi.NewRouter()

	idx := &Index{}
	router.Get("/", idx.ServeHTTP)

	nf := &NotFound{}
	router.NotFound(nf.ServeHTTP)

	return router
}
