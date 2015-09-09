package acceptance_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/franela/goreq"
)

var (
	TargetHost = os.Getenv("TARGET_HOST")
)

func TestMain(m *testing.M) {
	if TargetHost == "" {
		fmt.Println("Cannot run acceptance tests. Missing TARGET_HOST environment variable.")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestItWorks(t *testing.T) {
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
