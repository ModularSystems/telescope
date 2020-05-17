package scan

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// WPScan represents the data model of a WPScan
type WPScan struct {
	Error     error
	Every     time.Duration
	ScanName  string
	Stdout    string
	Stderr    string
	Time      string // Used to invoke after a certain time
	Timestamp time.Time
	Trigger   string
	URI       string
}

// GetName returns a string for easier lookups in maps
func (w *WPScan) GetName() string {
	return w.ScanName
}

// GetTimestamp returns the scan's timestamp for comparisons
func (w *WPScan) GetTimestamp() time.Time {
	return w.Timestamp
}

// IsEligible determines scan eligibility based on the current time
func (w *WPScan) IsEligible(lastRun time.Time) bool {
	durationSinceLastRun := time.Now().Sub(lastRun)

	if durationSinceLastRun > w.Every {
		return true
	}
	return false
}

// Scan runs wpscan and saves stdout/stderr output
// Scan looks for WPVULNAPITOKEN in the environment variables, and if it's there, enables
// wpvulndb lookups, which is recommended
func (w *WPScan) Scan() bool {
	t := time.Now()
	w.Timestamp = t
	if w.URI == "" {
		w.Error = errors.New("Scan cannot run if the URI isn't set")
		return false
	}

	// compose args
	args := fmt.Sprintf("wpscan --url %s", w.URI)
	if os.Getenv("WPVULNDB_API_KEY") != "" {
		args += fmt.Sprintf(" --api-token $WPVULNAPITOKEN")
	}

	// get `wpscan` executable path
	wpscanBinary, err := exec.LookPath("wpscan")
	if err != nil {
		w.Error = err
		return false
	}

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	outWriter := bufio.NewWriter(outBuf)
	errWriter := bufio.NewWriter(errBuf)

	splitArgs := strings.Split(args, " ")
	cmd := &exec.Cmd{
		Path:   wpscanBinary,
		Args:   splitArgs,
		Stdout: outWriter,
		Stderr: errWriter,
	}

	w.Error = cmd.Run()
	w.Stdout = outBuf.String()
	w.Stderr = errBuf.String()
	if w.Error != nil {
		return false
	}
	return true
}

// HTMLScan represents the data model of an HTML Scan
type HTMLScan struct {
	Error     error
	Every     time.Duration
	HTML      string
	ScanName  string
	Time      string
	Timestamp time.Time
	URI       string
}

// GetName returns a name for the scanner to be used in maps for easier lookups
func (h *HTMLScan) GetName() string {
	return h.ScanName
}

// GetTimestamp returns the scan's timestamp for comparisons
func (h *HTMLScan) GetTimestamp() time.Time {
	return h.Timestamp
}

// IsEligible determines scan eligibility based on the current time
func (h *HTMLScan) IsEligible(lastRun time.Time) bool {
	durationSinceLastRun := time.Now().Sub(lastRun)

	if durationSinceLastRun > h.Every {
		return true
	}
	return false
}

// Scan stores the HTML from a uri, unless there is an error in which case the error is saved
func (h *HTMLScan) Scan() bool {
	t := time.Now()
	h.Timestamp = t

	res, err := http.Get(h.URI)
	if err != nil {
		h.Error = err
		return false
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		h.Error = err
		return false
	}
	h.HTML = string(data)
	return true
}
