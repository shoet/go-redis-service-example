package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	srv := &http.Server{
		Handler: mux,
	}
	s := &Server{
		srv: srv,
		l:   l,
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed run serve: %w", err)
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}
	fmt.Println("shutdown server")

	return eg.Wait()
}
