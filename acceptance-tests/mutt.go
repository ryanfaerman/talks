package acceptance_test

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/franela/goreq"
	"github.com/nbio/st"
)

var (
	TargetHost = os.Getenv("TARGET_HOST")
)

func acceptanceTest(t *testing.T) {
	if testing.Short() {
		t.Skip("This is an acceptance test")
	}
	if TargetHost == "" {
		t.Fatal("Cannot run acceptance tests. Missing TARGET_HOST environment variable.")
	}

}

func TestItWorks(t *testing.T) {
	acceptanceTest(t) // HL

	res, err := goreq.Request{
		Uri:     TargetHost + "/health-check",
		Timeout: 1 * time.Second,
	}.Do()
	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()

	st.Assert(t, err, nil)
	st.Expect(t, res.StatusCode, http.StatusOK)

	body, err := res.Body.ToString()
	st.Assert(t, err, nil)
	st.Expect(t, strings.TrimSpace(body), "ok")
}
