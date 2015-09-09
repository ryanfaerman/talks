package acceptance_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/franela/goreq"
)

var (
	TargetHost = os.Getenv("TARGET_HOST")
)

func acceptanceTest(t *testing.T) {
	if TargetHost == "" {
		t.Fatal("Cannot run acceptance tests. Missing TARGET_HOST environment variable.")
	}
	if testing.Short() {
		t.Skip("This is an acceptance test")
	}
}

func TestItWorks(t *testing.T) {
	acceptanceTest(t)

	res, err := goreq.Request{
		Uri:     TargetHost,
		Timeout: 1 * time.Second,
	}.Do()
	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		t.Fatalf("Unexpected Request error: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected to receive %d but got %d instead", http.StatusOK, res.StatusCode)
	}
}
