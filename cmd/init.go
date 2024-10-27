package cmd

import (
	"context"
	"fmt"
	"github.com/nigel-campbell/pubsub/pubsub"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the PubSub environment",
	Long: `Sets up the necessary database tables and prepares the PubSub environment for operation.

The "init" command creates the required database schema, including tables for topics, subscriptions, and messages, 
in the SQLite database. Run this command once before adding topics, subscriptions, or messages to ensure the database 
is correctly set up.

Examples:
  pubsub init   # Sets up the database and required tables
`,
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
