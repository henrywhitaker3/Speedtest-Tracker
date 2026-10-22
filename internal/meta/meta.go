// Package meta
package meta

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Ready = false
)

type Server struct {
	server *http.Server
}

func New() *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if !Ready {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})

	return &Server{
		server: &http.Server{
			Addr:    ":8082",
			Handler: mux,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error)

	go func() {
		slog := slog.With("component", "meta")
		slog.Info("starting meta server", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}
	}()

	select {
	case <-ctx.Done():
		cctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		slog.Info("stopping meta server")
		return s.server.Shutdown(cctx)
	case err := <-errCh:
		return err
	}
}
