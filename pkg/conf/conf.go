package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ScanConfiguration represents the data model for a scanner's configuration.
type ScanConfiguration struct {
	Alerts []string          `yaml:"alerts"` // A list of alert's names to check
	Config map[string]string `yaml:"config"` // An optional config mapping
	Name   string            `yaml:"name"`   // An identifier for the Scan
	Type   string            `yaml:"type"`   // Triggers the behavior of the scan
	URIs   []string          `yaml:"uris"`   // The network endpoint(s) we want to trigger the scan for
}

// AlertConfiguration represents the data model from an alert configuration. This should map closely to the alert pkg object
type AlertConfiguration struct {
	Attribute string            `yaml:"attribute"` // The attribute we want to perform a regex against
	Config    map[string]string `yaml:"config"`    // User defined configuration options
	Name      string            `yaml:"name"`      // An identifier for the Alert
	Regex     string            `yaml:"regex"`     // Regex pattern matching to be applied against attribute
	Type      string            `yaml:"type"`      // The type of alerting to be performed (email, text, etc)
	URIs      []string          `yaml:"uris"`      // The network endpoint(s) relevant to the alert
}

// Config represents the data model for the configuration file loaded into telescope
type Config struct {
	Alerts []AlertConfiguration
	Scans  []ScanConfiguration
}

// Load unmarshals a file into our configuration struct
func (c *Config) Load(filepath string) error {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	return nil
}
