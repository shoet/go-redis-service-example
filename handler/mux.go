package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/go-redis-service-example/config"
	"github.com/shoet/go-redis-service-example/service"
	"github.com/shoet/go-redis-service-example/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, error) {
	router := chi.NewRouter()

	validator := validator.New()
	kvs, err := store.NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	idx := &Index{}
	router.Get("/", idx.ServeHTTP)

	l := &Login{
		Service: &service.LoginService{
			Store: kvs,
		},
		Validator: validator,
	}
	router.Post("/login", l.ServeHTTP)

	nf := &NotFound{}
	router.NotFound(nf.ServeHTTP)

	return router, nil
}
