package service

import (
	"context"
	"fmt"
	"log"

	"github.com/shoet/go-redis-service-example/store"
)

type LoginService struct {
	Store store.KVStore
}

func (ls *LoginService) Login(ctx context.Context, username string) error {
	log.Println("login service")
	v, err := ls.Store.Get(ctx, username)
	log.Printf("v: %s\n", v)
	if err != nil {
		// return fmt.Errorf("not found user: %v", err)
		fmt.Printf("kvs error: %v", err)
		fmt.Printf("user not found with login %s\n", username)
		ls.Store.Set(ctx, username, "login")
	}
	if v == "" {
		fmt.Printf("user not found with login %s\n", username)
		ls.Store.Set(ctx, username, "login")
	}
	return nil
}
