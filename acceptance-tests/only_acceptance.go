package acceptance_test

import (
	"fmt"
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

	st.Assert(t, err, nil)
	st.Expect(t, res.StatusCode, http.StatusOK)

	body, err := res.Body.ToString()
	st.Assert(t, err, nil)
	st.Expect(t, strings.TrimSpace(body), "ok")
}
