package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestLogin(t *testing.T) {
	mock := &LoginServiceMock{}
	mock.LoginFunc = func(ctx context.Context, username string) error {
		return nil
	}
	v := validator.New()
	sut := Login{
		Service:   mock,
		Validator: v,
	}

	b := []byte(`{"username": "test"}`)
	r := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
