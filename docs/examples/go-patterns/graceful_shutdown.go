// Package patterns demonstrates graceful shutdown for a Go service
// that runs multiple components (HTTP server, worker, etc.).
//
// Key rules from CLAUDE.md:
// - Use contexts for timeout handling and cancellation
// - Use errgroup for concurrent operations
package patterns

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const shutdownTimeout = 15 * time.Second

// RunService starts an HTTP server and a background worker,
// and shuts both down gracefully on SIGTERM or SIGINT.
func RunService(httpAddr string) error {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv := &http.Server{Addr: httpAddr}

	g, ctx := errgroup.WithContext(ctx)

	// Start HTTP server.
	g.Go(func() error {
		log.WithField("addr", httpAddr).Info("starting HTTP server")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	// Graceful shutdown goroutine.
	g.Go(func() error {
		<-ctx.Done()
		log.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}
		log.Info("http server stopped")
		return nil
	})

	return g.Wait()
}
