package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

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
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed run serve: %w", err)
		}
		return nil
	})

	return eg.Wait()
}
