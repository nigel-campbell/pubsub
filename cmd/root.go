/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// pubsubCmd represents the root command for managing the Pub/Sub emulator
var rootCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Manage topics, subscriptions, and messages in the Pub/Sub emulator",
	Long: `A command-line tool for managing a simulated Pub/Sub service.
This tool allows you to:
  - Create and manage topics and subscriptions
  - Publish and list messages
  - Acknowledge or nack messages
  - Initialize and clean up resources
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pubsub.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
