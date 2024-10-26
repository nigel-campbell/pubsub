package cmd

import (
	"context"
	"fmt"
	"log"
	"pubsub/pubsub"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// Use a default timeout for database operations
const dbTimeout = 5 * time.Second

// listCmd represents the base "list" command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List topics, subscriptions, or messages",
	Long:  "Use this command to list topics, subscriptions for a given topic, or messages for a specific topic and subscription.",
}

// listTopicsCmd lists all topics
var listTopicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List all topics",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Failed to initialize service: %v", err)
		}

		topics, err := svc.ListTopics(ctx)
		if err != nil {
			log.Fatalf("Error retrieving topics: %v", err)
		}

		fmt.Println("Topics:")
		for _, topic := range topics {
			fmt.Printf("- ID: %d, Name: %s, Metadata: %s\n", topic.ID, topic.Name, topic.Metadata)
		}
	},
}

// listSubscriptionsCmd lists subscriptions for a given topic
var listSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions [TOPIC_ID]",
	Short: "List all subscriptions for a specific topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Failed to initialize service: %v", err)
		}

		topicId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Error converting topic ID to integer: %v", err)
		}

		subscriptions, err := svc.ListSubscriptions(ctx, topicId)
		if err != nil {
			log.Fatalf("Error retrieving subscriptions for topic %d: %v", topicId, err)
		}
		fmt.Printf("Subscriptions for topic %d:\n", topicId)
		for _, sub := range subscriptions {
			fmt.Printf("- ID: %d, SubscriberID: %s\n", sub.ID, sub.SubscriberID)
		}
	},
}

// listMessagesCmd lists messages for a given topic and subscription
var listMessagesCmd = &cobra.Command{
	Use:   "messages [SUBSCRIPTION_ID]",
	Short: "List all messages for a specific subscription",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Failed to initialize service: %v", err)
		}

		subscriptionId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Error converting topic ID to integer: %v", err)
		}

		messages, err := svc.GetMessages(ctx, subscriptionId)
		if err != nil {
			log.Fatalf("Error retrieving messages for subscription %s: %v", subscriptionId, err)
		}

		if len(messages) > 0 {
			fmt.Printf("Messages for subscription %d:\n", subscriptionId)
			for _, msg := range messages {
				fmt.Println(msg)
			}
		} else {
			fmt.Println("No messages found for subscription", subscriptionId)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listTopicsCmd)
	listCmd.AddCommand(listSubscriptionsCmd)
	listCmd.AddCommand(listMessagesCmd)
}
