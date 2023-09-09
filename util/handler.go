package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/shoet/go-redis-service-example/errorutil"
	"github.com/shoet/go-redis-service-example/store"
	"golang.org/x/net/context"
)

func RespondJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageInternalServerError}
		err := json.NewEncoder(w).Encode(errResp)
		if err != nil {
			log.Printf("failed to encode error response: %v", err)
		}
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonBody)
	return
}

type JWT struct {
	JwtSecret string
	KVStore   store.KVStore
}

func (j *JWT) GenerateJWT(ctx context.Context, username string) (string, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Subject("auth_token").
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(time.Minute*30)).
		Claim("username", username).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build jwt: %w", err)
	}
	secrets := []byte(j.JwtSecret)
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, secrets))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return string(signed), nil
}

func (j *JWT) ValidateJWT(token string) (jwt.Token, error) {
	parsed, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, []byte(j.JwtSecret)), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}
	return parsed, nil

}

func (j *JWT) AuthGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tokenが無い、redisでのセッション切れはguardする
		token, err := r.Cookie("auth-token")
		if err != nil {
			if err == http.ErrNoCookie {
				errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageUnauthorized}
				RespondJSON(w, http.StatusUnauthorized, errResp)
				return
			}
			errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageInternalServerError}
			RespondJSON(w, http.StatusInternalServerError, errResp)
			return
		}
		t, err := j.ValidateJWT(token.Value)
		if err != nil {
			errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageInternalServerError}
			RespondJSON(w, http.StatusInternalServerError, errResp)
			return
		}
		username, exists := t.Get("username")
		if exists == false {
			errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageUnauthorized}
			RespondJSON(w, http.StatusUnauthorized, errResp)
			return
		}
		ret, err := j.KVStore.Get(r.Context(), username.(string))
		if err != nil {
			errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageUnauthorized}
			RespondJSON(w, http.StatusUnauthorized, errResp)
			return
		}
		fmt.Println("ret")
		fmt.Println(ret)
		if ret == "" {
			w.Header().Set("Set-Cookie", "auth-token=; Max-Age=0")
			errResp := errorutil.ErrorResponse{Message: errorutil.ErrorMessageUnauthorized}
			RespondJSON(w, http.StatusUnauthorized, errResp)
			return
		}
		next.ServeHTTP(w, r)
	})
}
