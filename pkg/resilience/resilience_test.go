package resilience

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestProcessMessage(t *testing.T) {
	msg := []string{"test message for processing", "a", "b", "c", "d", "e"}

	for _, m := range msg {
		err := ProcessMessage(m)
		// Since the function has a random failure, we can't test for a single outcome.
		// Instead, we check that the outcome is one of the two possibilities: success (nil error)

		// or the known simulated failure.
		if err != nil {
			expectedErr := "simulated processing failed"
			t.Logf("The err is %v", expectedErr)
			if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("ProcessMessage() returned an unexpected error. got: %v, want: %v", err, expectedErr)
			}
		}
	}
	// If err is nil, the test passes, as it's the success case.
}

// withTestContext is a helper for testing to inject a context.
func withTestContext(ctx context.Context) RetryOption {
	return func(args *retryOptions) {
		args.ctx = ctx
	}
}
func TestRetry(t *testing.T) {
	t.Run("SuccessOnFirstTry", func(t *testing.T) {
		var attempts int
		op := func() error {
			attempts++
			return nil
		}
		err := Retry(op, WithMaxTries(3))
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if attempts != 1 {
			t.Errorf("Expected 1 attempt, but got %d", attempts)
		}
	})

	t.Run("SuccessAfterRetries", func(t *testing.T) {
		attempts := 0
		op := func() error {
			attempts++
			if attempts < 3 {
				return errors.New("transient error")
			}
			return nil
		}
		err := Retry(op, WithMaxTries(5))
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, but got %d", attempts)
		}
	})

	t.Run("FailureAfterMaxRetries", func(t *testing.T) {
		attempts := 0
		op := func() error {
			attempts++
			return errors.New("persistent error")
		}
		err := Retry(op, WithMaxTries(3))
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if attempts != 2 {
			t.Errorf("Expected 2 attempts, but got %d", attempts)
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		op := func() error {
			cancel()
			return errors.New("some error")
		}

		// This will panic without a fix in retry.go
		err := Retry(op, withTestContext(ctx))

		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, but got %v", err)
		}
	})

	t.Run("MaxElapsedTimeExceeded", func(t *testing.T) {
		var attempts int
		op := func() error {
			attempts++
			time.Sleep(50 * time.Millisecond)
			return errors.New("persistent error")
		}
		err := Retry(op, WithMaxElapsedTime(100*time.Millisecond), WithMaxTries(10))
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if attempts > 2 {
			t.Errorf("Expected at most 2 attempts, but got %d", attempts)
		}
	})
}

func TestTimeout(t *testing.T) {
	t.Run("execution succeeded", func(t *testing.T) {
		var counter int
		err := WithTimeout(func() error {
			counter++
			return nil
		}, 1*time.Second)
		if err != nil {
			t.Errorf("\"1st\" test failed, expected %v, got %v", nil, err)
		}
		if counter != 1 {
			t.Errorf("\"1st\" test failed, Expected counter == 1, but got %d", counter)
		}
	})
	t.Run("execution timeout", func(t *testing.T) {
		var counter int
		err := WithTimeout(func() error {
			time.Sleep(2 * time.Second)
			counter++
			return nil
		}, 1*time.Second)
		if err != ErrTimeout {
			t.Errorf("\"2nd\" test failed, expected: %v, got %v", ErrTimeout, err)
		}
		if counter != 0 {
			t.Errorf("\"2nd\" test failed, Expected counter == 1, but got %d", counter)
		}
	})
}
