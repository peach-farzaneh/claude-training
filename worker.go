package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	defaultMaxRetries  = 3
	defaultBaseBackoff = 500 * time.Millisecond
	defaultTimeout     = 30 * time.Second
)

// Job represents a unit of work pulled from the Redis queue.
type Job struct {
	ID      string
	Payload string
}

// Worker reads jobs from a Redis list and processes them with retry logic.
type Worker struct {
	client     *redis.Client
	queue      string
	maxRetries int
	logger     *log.Entry
}

// NewWorker creates a Worker connected to the given Redis address and queue name.
func NewWorker(redisAddr, queue string) *Worker {
	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	return &Worker{
		client:     client,
		queue:      queue,
		maxRetries: defaultMaxRetries,
		logger:     log.WithField("component", "worker"),
	}
}

// Run starts the worker loop, polling the Redis queue until ctx is cancelled.
// It uses errgroup to manage concurrent processing goroutines.
func (w *Worker) Run(ctx context.Context, concurrency int) error {
	w.logger.WithField("concurrency", concurrency).Info("starting worker")

	g, ctx := errgroup.WithContext(ctx)
	for i := range concurrency {
		id := i
		g.Go(func() error {
			return w.poll(ctx, id)
		})
	}

	err := g.Wait()
	w.logger.Info("worker stopped")
	return err
}

func (w *Worker) poll(ctx context.Context, id int) error {
	logger := w.logger.WithField("goroutine", id)
	for {
		select {
		case <-ctx.Done():
			logger.Info("shutting down")
			return nil
		default:
		}

		job, err := w.dequeue(ctx)
		if err != nil {
			logger.WithError(err).Warn("dequeue failed, backing off")
			if sleepCtx(ctx, defaultBaseBackoff) != nil {
				return nil
			}
			continue
		}
		if job == nil {
			if sleepCtx(ctx, defaultBaseBackoff) != nil {
				return nil
			}
			continue
		}

		if err := w.processWithRetry(ctx, job); err != nil {
			logger.WithError(err).WithField("job_id", job.ID).Error("job failed after retries")
		}
	}
}

// dequeue pops a job from the Redis queue, returning nil when the queue is empty.
func (w *Worker) dequeue(ctx context.Context) (*Job, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	result, err := w.client.BLPop(ctx, defaultBaseBackoff, w.queue).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dequeueing from %s: %w", w.queue, err)
	}

	return &Job{
		ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
		Payload: result[1],
	}, nil
}

// processWithRetry attempts to process a job up to maxRetries times using
// exponential backoff between attempts.
func (w *Worker) processWithRetry(ctx context.Context, job *Job) error {
	var lastErr error
	for attempt := range w.maxRetries {
		logger := w.logger.WithFields(log.Fields{
			"job_id":  job.ID,
			"attempt": attempt + 1,
		})

		lastErr = ProcessJob(ctx, job)
		if lastErr == nil {
			logger.Info("job completed")
			return nil
		}

		logger.WithError(lastErr).Warn("job attempt failed")

		backoff := time.Duration(math.Pow(2, float64(attempt))) * defaultBaseBackoff
		if sleepCtx(ctx, backoff) != nil {
			return fmt.Errorf("context cancelled during retry backoff: %w", ctx.Err())
		}
	}
	return fmt.Errorf("job %s failed after %d attempts: %w", job.ID, w.maxRetries, lastErr)
}

// ProcessJob handles the business logic for a single job.
func ProcessJob(ctx context.Context, job *Job) error {
	log.WithFields(log.Fields{
		"job_id":  job.ID,
		"payload": job.Payload,
	}).Info("processing job")

	select {
	case <-ctx.Done():
		return fmt.Errorf("job %s cancelled: %w", job.ID, ctx.Err())
	default:
	}

	if job.Payload == "" {
		return fmt.Errorf("job %s has empty payload: %w", job.ID, ErrEmptyPayload)
	}

	return nil
}

// ErrEmptyPayload is returned when a job has no payload data.
var ErrEmptyPayload = fmt.Errorf("empty payload")

// Close shuts down the Redis connection.
func (w *Worker) Close() error {
	if err := w.client.Close(); err != nil {
		return fmt.Errorf("closing redis client: %w", err)
	}
	return nil
}

// sleepCtx pauses for the given duration, returning early if ctx is cancelled.
func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
