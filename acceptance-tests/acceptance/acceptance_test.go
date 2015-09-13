package main_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"camlistore.org/pkg/netutil"

	"github.com/franela/goreq"
	"github.com/nbio/st"
)

const (
	ApplicationHost = "127.0.0.1:8080"
	StartupTimeout  = 5
)

var (
	ApplicationName    = findApplicationName()
	ApplicationPIDFile = ApplicationName + ".pid"
)

func findApplicationName() string {
	pwd, err := os.Getwd()
	if err != nil {
		abortLog(err)
	}
	return filepath.Base(pwd)
}

func TestMain(m *testing.M) {
	flag.Parse()
	killPrevious()

	if err := build(); err != nil {
		abortLog(err)
	}
	process, err := run()
	if err != nil {
		abortLog(err)
	}
	writePidFile(process.Pid)

	if err := waitForReachable(StartupTimeout * time.Second); err != nil {
		verboseLog("Application is unreachable after", StartupTimeout, "seconds")
		kill(process)
		abortLog("Cannot start tests")
	}

	verboseLog("Starting test suite")
	suite := m.Run()
	kill(process)
	os.Exit(suite)
}

func build() error {
	verboseLog("Building application", ApplicationName)

	cmd := exec.Command("go", "build")

	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v%v", stderr.String(), err)
	}

	return nil
}

func run() (*os.Process, error) {
	verboseLog("Starting", ApplicationName, "as a child process")

	cmd := exec.Command("./" + ApplicationName)

	if testing.Verbose() { // HL
		cmd.Stdout = os.Stdout // HL
		cmd.Stderr = os.Stderr // HL
	} // HL

	if err := cmd.Start(); err != nil { // HL
		log.Fatal(err) // HL
	} // HL

	verboseLog("Child process started with PID:", cmd.Process.Pid)
	return cmd.Process, nil
}

func writePidFile(pid int) {
	file, err := os.Create(ApplicationPIDFile)
	defer file.Close()

	if err != nil {
		abortLog(err)
	}

	fmt.Fprint(file, pid)
}

func kill(p *os.Process) {
	verboseLog("Killing child process")
	if err := p.Kill(); err != nil {
		traceLog("Could not kill", ApplicationName, "child is still running with PID:", p.Pid)
	}

	verboseLog("Removing test artifacts")
	os.Remove(ApplicationPIDFile)
	os.Remove(ApplicationName)
}

func killPrevious() {
	content, err := ioutil.ReadFile(ApplicationPIDFile)
	if err == nil {
		verboseLog("Detected previous running instance with PID:", string(content))

		pid, _ := strconv.Atoi(string(content))
		p, _ := os.FindProcess(pid)
		kill(p)
	}
}

func waitForReachable(timeout time.Duration) error {
	verboseLog("Waiting up to", StartupTimeout, "seconds for", ApplicationName, "to be reachable at", ApplicationHost)

	if err := netutil.AwaitReachable(ApplicationHost, timeout); err != nil { // HL
		return err // HL
	} // HL

	verboseLog("Application Host is reachable, resuming")
	return nil
}

func verboseLog(msg ...interface{}) {
	if testing.Verbose() {
		traceLog(msg...)
	}
}

func traceLog(msg ...interface{}) {
	fmt.Println(append([]interface{}{"-->"}, msg...)...)
}

func abortLog(err interface{}) {
	traceLog(err)
	os.Exit(1)
}

func acceptanceTest(t *testing.T) {
	if testing.Short() {
		t.Skip("This is an acceptance test")
	}
}

func TestItWorks(t *testing.T) {
	acceptanceTest(t)

	res, err := goreq.Request{
		Uri:     "http://" + ApplicationHost, // HL
		Timeout: 1 * time.Second,
	}.Do()
	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()

	st.Assert(t, err, nil)
	st.Expect(t, res.StatusCode, http.StatusOK)
}
