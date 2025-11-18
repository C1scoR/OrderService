package resilience

import (
	"strings"
	"testing"
)

func TestProcessMessage(t *testing.T) {
	msg := "test message for processing"
	err := ProcessMessage(msg)

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
	// If err is nil, the test passes, as it's the success case.
}
