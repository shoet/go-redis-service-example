package handler

import (
	"context"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . LoginService
type LoginService interface {
	Login(ctx context.Context, username string) error
}
