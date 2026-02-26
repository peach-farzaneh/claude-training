// Package patterns demonstrates exponential backoff retry logic.
//
// Key rule from CLAUDE.md:
// - Never use time.Sleep() for retries — use exponential backoff
// - Use context.Context in all functions that do I/O
package patterns

import (
	"context"
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
)

// RetryConfig controls the retry behavior.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns sensible defaults for most operations.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    30 * time.Second,
	}
}

// WithRetry executes fn up to MaxAttempts times with exponential backoff.
// It respects context cancellation and never uses time.Sleep.
func WithRetry(ctx context.Context, cfg RetryConfig, operation string, fn func(ctx context.Context) error) error {
	var lastErr error

	for attempt := range cfg.MaxAttempts {
		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}

		log.WithFields(log.Fields{
			"operation": operation,
			"attempt":   attempt + 1,
			"error":     lastErr,
		}).Warn("operation failed, retrying")

		if attempt == cfg.MaxAttempts-1 {
			break
		}

		delay := time.Duration(math.Pow(2, float64(attempt))) * cfg.BaseDelay
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("%s cancelled during retry: %w", operation, ctx.Err())
		case <-timer.C:
		}
	}

	return fmt.Errorf("%s failed after %d attempts: %w",
		operation, cfg.MaxAttempts, lastErr)
}
