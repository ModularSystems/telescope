package cmd

import (
	"fmt"
	"log"
	"os"

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

			logger := log.New(os.Stdout, "", log.LstdFlags)
			if debug {
				logger.Printf("✔️ Configuration file loaded from %s\n", configFile)
				if os.Getenv("SENDGRID_API_KEY") != "" && os.Getenv("SENDGRID_SENDER_NAME") != "" && os.Getenv("SENDGRID_SENDER_EMAIL") != "" {
					logger.Printf("✔️ Sendgrid enabled\n")
				} else {
					logger.Printf("✖ Sendgrid disabled\n")
				}
				if os.Getenv("WPVULNDB_API_KEY") != "" {
					logger.Printf("✔️ WPVulnDB lookups enabled\n")
				} else {
					logger.Printf("✖ WPVulnDB lookups disabled\n")
				}
			}
			store := &daemon.InMemoryStore{
				CacheLength: 100,
			}
			daemon := &daemon.Daemon{
				Config:  config,
				Debug:   debug,
				Logger:  logger,
				Storage: store,
			}
			daemon.Load()
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
