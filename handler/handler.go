package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/go-redis-service-example/errorutil"
	"github.com/shoet/go-redis-service-example/util"
)

type Index struct{}

func (*Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AuthGuard(w, r)
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}
	util.RespondJSON(w, http.StatusOK, resp)
}

func AuthGuard(w http.ResponseWriter, r *http.Request) {
	// tokenが無い、redisでのセッション切れはguardする
	token, err := r.Cookie("auth-token")
	if err != nil {
		if err == http.ErrNoCookie {
			errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageUnauthorized}
			util.RespondJSON(w, http.StatusUnauthorized, errResp)
			return
		}
		log.Printf("failed to get cookie: %v", err)
		errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageInternalServerError}
		util.RespondJSON(w, http.StatusInternalServerError, errResp)
		return
	}
	// TODO: token 検証
	if !util.ValidateJWT(token.Value) {
		errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageUnauthorized}
		util.RespondJSON(w, http.StatusUnauthorized, errResp)
	}
	// TODO: redis 検索
	// TODO: cookie reset
}

type Login struct {
	Service   LoginService
	Validator *validator.Validate
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
	resp := struct {
		Name string `json:"name"`
	}{
		Name: body.Username,
	}
	expire := time.Now()
	expire = expire.AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "auth-token", Value: "test", Expires: expire}
	w.Header().Set("Set-Cookie", cookie.String())
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
