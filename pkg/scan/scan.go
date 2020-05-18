package scan

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Scanner sets up a common interface for a variety of scans
type Scanner interface {
	GetError() error // should return whatever error output is logical to pass to the user
	GetName() string
	GetOutput() string // should return whatever is logical for the scan to output
	GetURI() string    // should return whatever is logical for the scan to output
	GetTimestamp() time.Time
	IsEligible(time.Time) bool // Used to determine if a scan is eligible to run
	Scan() bool                // Scan returns true when a scan is performed
}

// WPScan represents the data model of a WPScan
type WPScan struct {
	Debug     bool
	Error     error
	Every     time.Duration
	Logger    *log.Logger
	ScanName  string
	Stdout    string
	Stderr    string
	Time      string // Used to invoke after a certain time
	Timestamp time.Time
	URI       string
}

// GetError returns any error output from the result of the scan
func (w *WPScan) GetError() error {
	return fmt.Errorf("stderr: %s\tlast error: %s\t", w.Stderr, w.Error.Error())
}

// GetName returns a string for easier lookups in maps
func (w *WPScan) GetName() string {
	return w.ScanName
}

// GetOutput returns whatever output is expected of the scan
func (w *WPScan) GetOutput() string {
	return w.Stdout
}

// GetURI returns the uri for identification
func (w *WPScan) GetURI() string {
	return w.URI
}

// GetTimestamp returns the scan's timestamp for comparisons
func (w *WPScan) GetTimestamp() time.Time {
	return w.Timestamp
}

// IsEligible determines scan eligibility based on the current time
func (w *WPScan) IsEligible(lastRun time.Time) bool {
	if w.Debug {
		w.Logger.Printf("Evaluating %s\tLast run: %v\t", w.ScanName, lastRun)
	}
	durationSinceLastRun := time.Now().Sub(lastRun)
	if w.Debug {
		w.Logger.Printf("Duration since last run: %v\tScan every: %v\n", durationSinceLastRun, w.Every)
	}
	if durationSinceLastRun > w.Every {
		return true
	}

	return false
}

// Scan runs wpscan and saves stdout/stderr output
// Scan looks for WPVULNDB_API_KEY in the environment variables, and if it's there, enables
// wpvulndb lookups, which is recommended
func (w *WPScan) Scan() bool {
	t := time.Now()
	w.Timestamp = t
	if w.URI == "" {
		w.Error = errors.New("Scan cannot run if the URI isn't set")
		return false
	}

	// compose args
	if w.Debug {
		if os.Getenv("WPVULNDB_API_KEY") != "" {
			w.Logger.Printf("✔️\t WPVulnDB lookups enabled \n")
		} else {
			w.Logger.Printf("✖\t WPVulnDB lookups disabled\n")
		}
	}
	// get `wpscan` executable path
	wpscanBinary, err := exec.LookPath("wpscan")
	if err != nil {
		w.Error = err
		return false
	}
	args := fmt.Sprintf("%s --url %s", wpscanBinary, w.URI)
	// TODO - Get this tp pick up the environment variable.
	if os.Getenv("WPVULNDB_API_KEY") != "" {
		// args += fmt.Sprintf(" --api-token $WPVULNDB_API_KEY")
		args += fmt.Sprintf(" --api-token %s", os.Getenv("WPVULNDB_API_KEY"))
	}

	errBuf := &bytes.Buffer{}
	errWriter := bufio.NewWriter(errBuf)
	outBuf := &bytes.Buffer{}
	outWriter := bufio.NewWriter(outBuf)

	splitArgs := strings.Split(args, " ")

	cmd := exec.Command(wpscanBinary, splitArgs...)
	cmd.Stderr = errWriter
	cmd.Stdout = outWriter

	w.Error = cmd.Run()
	w.Stderr = errBuf.String()
	w.Stdout = outBuf.String()
	if w.Error != nil || w.Stderr != "" {
		return false
	}
	return true
}

// HTMLScan represents the data model of an HTML Scan
type HTMLScan struct {
	Debug     bool
	Error     error
	Every     time.Duration
	Logger    *log.Logger
	HTML      string
	ScanName  string
	Time      string
	Timestamp time.Time
	URI       string
}

// GetError returns any error output from the result of the scan
func (h *HTMLScan) GetError() error {
	return h.Error
}

// GetName returns a name for the scanner to be used in maps for easier lookups
func (h *HTMLScan) GetName() string {
	return h.ScanName
}

// GetOutput returns whatever output is expected of the scan
func (h *HTMLScan) GetOutput() string {
	return h.HTML
}

// GetURI returns the uri for identification
func (h *HTMLScan) GetURI() string {
	return h.URI
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
