package cmd

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/modularsystems/telescope"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	configFile     string
	debug 	bool

	rootCmd = &cobra.Command{
		Use:   "telescope",
		Short: "A website monitor for changes, outages, and vulnerabilities",
		Long: `A website monitor for changes, outages, and vulnerabilities`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is /etc/telescope/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", true, "enable debug logging")
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		configFile = "/etc/telescope/config.yaml"
	}

	config := &Config{}
	err := config.Load(configFile)
	if err != nil {
		fmt.Printf("Failed to load configuration file at %s\t%s", configFile, err.Error())
		os.Exit(1)
	}
}