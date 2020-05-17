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
	Name      string
	Stdout    string
	Stderr    string
	Time      string // Used to invoke after a certain time
	Timestamp time.Time
	URI       string
}

// Scan runs wpscan and saves stdout/stderr output
// Scan looks for WPVULNAPITOKEN in the environment variables, and if it's there, enables
// wpvulndb lookups, which is recommended
func (w *WPScan) Scan() {
	t := time.Now()
	w.Timestamp = t
	if w.URI == "" {
		w.Error = errors.New("Scan cannot run if the URI isn't set")
		return
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
		return
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
		return
	}
}

// HTMLScan represents the data model of an HTML Scan
type HTMLScan struct {
	Error     error
	HTML      string
	Name      string
	Time      string
	Timestamp time.Time
	URI       string
}

// Scan stores the HTML from a uri, unless there is an error in which case the error is saved
func (h *HTMLScan) Scan() {
	t := time.Now()
	h.Timestamp = t

	res, err := http.Get(h.URI)
	if err != nil {
		h.Error = err
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		h.Error = err
		return
	}
	h.HTML = string(data)
}
