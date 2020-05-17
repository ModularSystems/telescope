package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Alert represents the data model from an alert configuration. This should map closely to the alert pkg object
type Alert struct {
	Name      string   `yaml:"name"`
	Scanner   string   `yaml:"scanner"`
	URIs      []string `yaml:"uris"`
	Attribute string   `yaml:"attribute"`
	Regex     string   `yaml:"regex"`
}

// Config represents the data model for the configuration file loaded into telescope
type Config struct {
	Alerts []Alert
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
