package resilience

import (
	"context"
	"math"
	"time"
)

//Kek lol

const (
	DefaultMaxElapsedTime = 2 * time.Minute
	DefaultMaxRetries     = 25
)

type retryOptions struct {
	backoff        *BackOff
	maxRetries     uint
	maxElapsedTime time.Duration
	ctx            context.Context
}

type RetryOption func(*retryOptions)

type BackOff struct {
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration

	currentInterval time.Duration
}

func (b *BackOff) NextBackOff() time.Duration {
	if b.currentInterval == 0 {
		b.currentInterval = b.InitialInterval
		return b.currentInterval
	}

	next := float64(b.currentInterval) * b.Multiplier
	b.currentInterval = time.Duration(math.Min(next, float64(b.MaxInterval)))
	return b.currentInterval
}

const (
	DefaultInitialInterval = 100 * time.Millisecond
	DefaultMultiplier      = 1.5
	DefaultMaxInterval     = 5 * time.Second
)

func NewBackoff() *BackOff {
	return &BackOff{
		InitialInterval: DefaultInitialInterval,
		Multiplier:      DefaultMultiplier,
		MaxInterval:     DefaultMaxInterval,
	}
}

func WithMaxTries(maxRetries uint) RetryOption {
	return func(args *retryOptions) {
		args.maxRetries = maxRetries
	}
}

func WithMaxElapsedTime(maxElapsedTime time.Duration) RetryOption {
	return func(args *retryOptions) {
		args.maxElapsedTime = maxElapsedTime
	}
}

func Retry(operation func() error, opts ...RetryOption) error {
	args := &retryOptions{
		backoff:        NewBackoff(),
		maxRetries:     DefaultMaxRetries,
		maxElapsedTime: DefaultMaxElapsedTime,
	}

	for _, opt := range opts {
		opt(args)
	}
	var lastErr error
	starttime := time.Now()
	for retr := uint(1); retr < args.maxRetries; retr++ {
		select {
		case <-args.ctx.Done():
			return args.ctx.Err()
		default:
			if time.Since(starttime) > args.maxElapsedTime {
				return lastErr
			}
			lastErr = operation()
			if lastErr == nil {
				return nil
			}
			if retr == args.maxRetries {
				break
			}

			sleepDuration := args.backoff.NextBackOff()
			select {
			case <-time.After(sleepDuration):
			case <-args.ctx.Done():
				return args.ctx.Err()
			}
		}
	}
	return lastErr
}
