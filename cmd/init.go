/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"pubsub-cli/pubsub"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			fmt.Println("Error creating Pub/Sub service:", err)
			return
		}
		defer svc.Close()

		err = svc.Init(context.Background())
		if err != nil {
			fmt.Println("Error initializing Pub/Sub service:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
