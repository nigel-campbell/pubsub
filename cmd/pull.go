package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pubsub/pubsub"
	"strconv"
	"time"
)

var ackDeadlineDuration time.Duration

var pullCmd = &cobra.Command{
	Use:   "pull [SUBSCRIPTION_ID]",
	Short: "Pull unacknowledged messages from a subscription",
	Long: `Retrieve messages from a specified subscription that have not been acknowledged.
You can also set an acknowledgment deadline using the flag`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Error creating Pub/Sub service: %v", err)
		}
		defer svc.Close()

		subscriptionID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid subscription ID: %s", args[0])
		}

		ackDeadline := time.Now().Add(ackDeadlineDuration)

		messages, err := svc.PullMessages(ctx, subscriptionID, ackDeadline)
		if err != nil {
			log.Fatalf("Failed to pull messages for subscription %s: %v", subscriptionID, err)
		}

		if len(messages) == 0 {
			fmt.Println("No messages to pull.")
		} else {
			fmt.Printf("Pulled messages for subscription %s:\n", subscriptionID)
			for _, msg := range messages {
				fmt.Println(msg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().DurationVarP(&ackDeadlineDuration, "deadline", "d", time.Second*10, "Set the acknowledgment deadline for pulled messages (e.g., 1m, 2h)")
}
