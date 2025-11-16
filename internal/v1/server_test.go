package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned the wrong Status: got: %v, want: %v", rr.Code, http.StatusOK)
		return
	}

	expected := "ok!200"
	if rr.Body.String() != expected {
		t.Errorf("HealthCheck returned the wrong Body, got: %v, expected: %s", rr.Body.String(), expected)
	}
}
