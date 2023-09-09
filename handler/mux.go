package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/go-redis-service-example/config"
	"github.com/shoet/go-redis-service-example/service"
	"github.com/shoet/go-redis-service-example/store"
)

func NewMux(cfg *config.Config) http.Handler {
	router := chi.NewRouter()

	validator := validator.New()
	kvs := store.NewRedisClient(cfg)

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

	return router
}
