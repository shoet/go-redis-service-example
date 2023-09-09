package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/go-redis-service-example/errorutil"
	"github.com/shoet/go-redis-service-example/util"
)

type Index struct{}

func (*Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}
	util.RespondJSON(w, http.StatusOK, resp)
}

type Login struct {
	Service   LoginService
	Validator *validator.Validate
	JWT       *util.JWT
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Username string `json:"username" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.RespondJSON(w, http.StatusInternalServerError, errorutil.ErrorResponse{Message: err.Error()})
		return
	}
	if err := l.Validator.Struct(&body); err != nil {
		util.RespondJSON(w, http.StatusBadRequest, errorutil.ErrorResponse{Message: err.Error()})
		return
	}
	if err := l.Service.Login(ctx, body.Username); err != nil {
		util.RespondJSON(w, http.StatusInternalServerError, errorutil.ErrorResponse{Message: err.Error()})
		return
	}
	expire := time.Now()
	expire = expire.AddDate(0, 0, 1)
	token, err := l.JWT.GenerateJWT(ctx, body.Username)
	if err != nil {
		util.RespondJSON(
			w,
			http.StatusInternalServerError,
			errorutil.ErrorResponse{Message: err.Error()},
		)
		return
	}
	cookie := http.Cookie{Name: "auth-token", Value: token, Expires: expire}
	w.Header().Set("Set-Cookie", cookie.String())
	resp := struct {
		Name string `json:"name"`
	}{
		Name: body.Username,
	}
	util.RespondJSON(w, http.StatusOK, resp)
}

type NotFound struct{}

func (*NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string `json:"message"`
	}{
		Message: errorutil.ErrorMessageNotFound,
	}
	util.RespondJSON(w, http.StatusNotFound, resp)
}
