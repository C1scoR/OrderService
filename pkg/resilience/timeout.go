package resilience

import (
	"errors"
	"time"
)

var ErrTimeout = errors.New("timeout exeeded")

func WithTimeout(operation func() error, duration time.Duration) error {
	errch := make(chan error)

	go func() {
		errch <- operation()
	}()

	go func() {
		time.Sleep(duration)
		errch <- ErrTimeout
	}()
	return <-errch
}
