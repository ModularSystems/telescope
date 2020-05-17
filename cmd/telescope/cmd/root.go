package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/modularsystems/telescope/pkg/conf"
	"github.com/modularsystems/telescope/pkg/daemon"
	"github.com/spf13/cobra"
)

var (
	config *conf.Config

	// Used for flags.
	configFile string
	debug      bool

	rootCmd = &cobra.Command{
		Use:   "telescope",
		Short: "A website monitor for changes, outages, and vulnerabilities",
		Long:  `A website monitor for changes, outages, and vulnerabilities`,
		Run: func(cmd *cobra.Command, args []string) {
			if debug {
				fmt.Println("debug logging enabled")
			}

			logger := log.New(os.Stdout, fmt.Sprintf("%s: ", time.Now().Format(time.RFC3339)), log.LUTC)

			daemon := daemon.Daemon{
				Config: config,
				Debug:  debug,
				Logger: logger,
			}
			daemon.Start()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "/etc/telescope/config.yaml", "config file (default is /etc/telescope/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")

}

func initConfig() {
	config = &conf.Config{}
	err := config.Load(configFile)
	if err != nil {
		fmt.Printf("Failed to load configuration file at %s\t%s", configFile, err.Error())
		os.Exit(1)
	}
}
