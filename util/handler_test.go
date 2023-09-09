package util

import (
	"context"
	"testing"
)

func TestJWT(t *testing.T) {
	ctx := context.Background()
	username := "test"
	jwt := JWT{JwtSecret: "testsecret"}
	token, err := jwt.GenerateJWT(ctx, username)
	if err != nil {
		t.Errorf("failed to generate jwt: %v", err)
	}

	v, err := jwt.ValidateJWT(token)
	if err != nil {
		t.Errorf("failed to validate jwt: %v", err)
	}
	if !v {
		t.Errorf("jwt is invalid")
	}
}
