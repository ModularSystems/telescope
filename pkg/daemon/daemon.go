package daemon

import (
	"log"
	"os"
	"strings"

	"github.com/modularsystems/telescope/pkg/alert"
	"github.com/modularsystems/telescope/pkg/conf"
	"github.com/modularsystems/telescope/pkg/scan"
)

// Daemon is used to configure and run the telescope daemon
type Daemon struct {
	Alerts []alert.Alert
	Config *conf.Config
	Debug  bool
	Logger *log.Logger
	Scans  []scan.Scanner
}

// Start loads everything into memory and starts our daemon
func (d *Daemon) Start() {
	d.Logger.Printf("Configuring scans\t")
	if d.Debug {
		d.Logger.Printf("Loaded Configuration\n")
		d.Logger.Printf("Loaded Alerts: %s\n", d.Config.Alerts)
		d.Logger.Printf("Loaded Scans: %s\n", d.Config.Scans)
	}
	for _, i := range d.Config.Scans {
		if d.Debug {
			d.Logger.Printf("found %s\t", i.Name)
		}
		if i.Type == "WPScan" {
			for _, uri := range i.URIs {
				WPScan := &scan.WPScan{
					Name: i.Name,
					URI:  uri,
					Time: i.Time,
				}
				d.Scans = append(d.Scans, WPScan)
			}
		}
		if i.Type == "HTMLScan" {
			for _, uri := range i.URIs {
				HTMLScan := &scan.HTMLScan{
					Name: i.Name,
					URI:  uri,
					Time: i.Time,
				}
				d.Scans = append(d.Scans, HTMLScan)
			}
		}
	}

	d.Logger.Printf("Configuring alerts\t")
	for _, a := range d.Config.Alerts {
		var toName, toEmail string
		if d.Debug {
			d.Logger.Printf("found %s\t", a.Name)
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

			fromName := os.Getenv("SENDGRID_SENDER_NAME")
			fromEmail := os.Getenv("SENDGRID_SENDER_EMAIL")
			if fromName == "" || fromEmail == "" {
				d.Logger.Printf("Unable to send sendgrid emails, SENDGRID_SENDER_NAME or SENDGRID_SENDER_EMAIL is unset\n")
				return
			}
			// TODO: Check that the config is properly constructed here
			emailAlert := alert.NewEmailAlert(fromName, fromEmail, toName, toEmail, a.Config["subject"], a.Config["message"])
			d.Alerts = append(d.Alerts, emailAlert)
		}
	}
}
