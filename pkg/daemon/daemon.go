package daemon

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/modularsystems/telescope/pkg/alert"
	"github.com/modularsystems/telescope/pkg/conf"
	"github.com/modularsystems/telescope/pkg/scan"
)

// Daemon is used to configure and run the telescope daemon
type Daemon struct {
	Alerts  map[string][]alert.Alert
	Config  *conf.Config
	Debug   bool
	Logger  *log.Logger
	Scans   map[string][]scan.Scanner
	Storage Storage
}

// Load sets up the Daemon for the main loop
func (d *Daemon) Load() {
	d.Scans = make(map[string][]scan.Scanner)
	d.Alerts = make(map[string][]alert.Alert)
	for _, i := range d.Config.Scans {
		if d.Debug {
			d.Logger.Printf("Loading scan: %s\t", i.Name)
		}
		if i.Type == "WPScan" {
			duration, err := time.ParseDuration(i.Config["every"])
			if err != nil {
				d.Logger.Printf("Error: Could not parse duration \"%s\" for alert %s please check the every attribute for proper formatting.", i.Config["every"], i.Name)
				d.Logger.Printf("Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\"\n")
				d.Logger.Printf("Defaulting to 24h for %s\n", i.Name)
			}
			for _, uri := range i.URIs {
				WPScan := &scan.WPScan{
					Every:    duration,
					Time:     i.Config["time"],
					ScanName: i.Name,
					URI:      uri,
				}
				d.Scans[i.Name] = append(d.Scans[i.Name], WPScan)
			}
		}
		if i.Type == "HTMLScan" {
			duration, err := time.ParseDuration(i.Config["every"])
			if err != nil {
				d.Logger.Printf("Error: Could not parse duration \"%s\" for alert %s please check the every attribute for proper formatting.", i.Config["every"], i.Name)
				d.Logger.Printf("Defaulting to 1d for %s\n", i.Name)
			}
			for _, uri := range i.URIs {
				HTMLScan := &scan.HTMLScan{
					Every:    duration,
					ScanName: i.Name,
					URI:      uri,
					Time:     i.Config["time"],
				}
				d.Scans[i.Name] = append(d.Scans[i.Name], HTMLScan)
			}
		}
	}

	if d.Debug {
		d.Logger.Printf("Configuring alerts\t")
	}
	for _, a := range d.Config.Alerts {
		var toName, toEmail string
		if d.Debug {
			d.Logger.Printf("Loading alert: %s\t", a.Name)
		}
		if a.Type == "email" {
			sendToCommaDelim := strings.Split(a.Config["sendTo"], ",")
			for _, sendTo := range sendToCommaDelim {
				tmp := strings.Split(sendTo, " ")
				if len(tmp) == 3 {
					toName = strings.Join(tmp[0:1], " ")
					toEmail = tmp[2]
				} else {
					if d.Debug {
						d.Logger.Printf("Failed to parse sender for %s\tconfig[\"sendTo\"]: %s", a.Name, a.Config["sendTo"])
					}
				}
			}

			// Check if email is enabled. If we don't have a valid email configuration, continue and see if we can load
			// other alerts.
			fromName := os.Getenv("SENDGRID_SENDER_NAME")
			fromEmail := os.Getenv("SENDGRID_SENDER_EMAIL")
			if fromName == "" || fromEmail == "" {
				continue
			}
			// TODO: Check that the config is properly constructed here
			emailAlert := alert.NewEmailAlert(a.Name, fromName, fromEmail, toName, toEmail, a.Config["subject"], a.Config["message"])
			d.Alerts[a.Name] = append(d.Alerts[a.Name], emailAlert)
		}
	}
	if d.Debug {
		d.Logger.Printf("Daemon configuration loaded")
	}

}

// Start loads everything into memory and starts our daemon
func (d *Daemon) Start() {
	if d.Debug {
		d.Logger.Printf("Daemon started")
	}
	tick := time.Tick(1 * time.Minute)
	for {
		select {
		case <-tick:
			// Iterate through scans' keys, and loop through each set of associated scanners. Each scanner should determine if it should be executed
			for k, v := range d.Scans {
				if d.Debug {
					d.Logger.Printf("Evaluating: %s\n", k)
				}
				// Iterate through the slice of Scanners
				for _, s := range v {
					var lastRun time.Time
					if d.Storage.SizeOf(s.GetName()) > 0 {
						lastRun = d.Storage.Last(s).GetTimestamp()
					}
					if s.IsEligible(lastRun) {
						if d.Debug {
							d.Logger.Printf("%s eligible for run, scanning now\n", s.GetName())
						}
						s.Scan()
						d.Storage.Save(s)
					}
				}
			}
		}
	}
}
