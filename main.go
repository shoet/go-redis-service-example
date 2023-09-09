package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/shoet/go-redis-service-example/config"
	"github.com/shoet/go-redis-service-example/handler"
)

func run() error {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	fmt.Printf("listen on http://%s\n", l.Addr().String())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	mux := handler.NewMux(cfg)
	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
