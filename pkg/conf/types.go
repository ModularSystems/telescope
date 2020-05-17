package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Scan represents the data model for a scanner's configuration.
type Scan struct {
	Alert  []string          `yaml:"alert"`  // A list of alert's names to check
	Config map[string]string `yaml:"config"` // An optional config mapping
	Name   string            `yaml:"name"`
	Time   string            `yaml:"time"` // The ISO 8061 Time string for when we want the Scan to happen
	Type   string            `yaml:"type"`
	URIs   []string          `yaml:"uris"`
}

// Alert represents the data model from an alert configuration. This should map closely to the alert pkg object
type Alert struct {
	Attribute string            `yaml:"attribute"`
	Config    map[string]string `yaml:"config"`
	Name      string            `yaml:"name"`
	Regex     string            `yaml:"regex"`
	Type      string            `yaml:"type"`
	URIs      []string          `yaml:"uris"`
}

// Config represents the data model for the configuration file loaded into telescope
type Config struct {
	Alerts []Alert
	Scans  []Scan
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
