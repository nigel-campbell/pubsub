/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log"
	"pubsub-cli/pubsub"
	"strconv"
)

// Define variables to store flags (like -d for config or payload)
var configFile string
var messagePayload string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add topics, subscriptions, or messages",
	Long:  "Use this command to add new topics, subscriptions, or messages to the Pub/Sub system",
}

var addTopicCmd = &cobra.Command{
	Use:   "topic [TOPIC_ID]",
	Short: "Add a new topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topicID := args[0]
		fmt.Printf("Adding topic: %s\n", topicID)

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			fmt.Println("Error creating Pub/Sub service:", err)
			return
		}
		defer svc.Close()

		err = svc.CreateTopic(context.Background(), topicID, []byte{})
		if err != nil {
			fmt.Println("Error creating topic:", err)
			return
		}

		fmt.Println("Topic created successfully")
	},
}

var addSubscriptionCmd = &cobra.Command{
	Use:   "subscription [TOPIC_ID] [SUBSCRIPTION_ID]",
	Short: "Add a new subscription to a topic",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		topicId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid topic ID: %s", args[0])
		}
		subscriptionName := args[1]
		// Implement the logic for adding a subscription using the topicID, subscriptionName, and configFile
		fmt.Printf("Adding subscription: %s to topic: %d\n", subscriptionName, topicId)

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Error creating Pub/Sub service: %s", err)
		}
		defer svc.Close()

		err = svc.CreateSubscription(context.Background(), topicId, subscriptionName, []byte{})
		if err != nil {
			log.Fatalf("Error creating subscription: %s", err)
		}
		fmt.Println("Subscription created successfully")
	},
}

// addMessageCmd represents the "add message" command
var addMessageCmd = &cobra.Command{
	Use:   "message [TOPIC_ID]",
	Short: "Add a message to a topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic for adding a message using the topicID and messagePayload
		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Error creating Pub/Sub service: %s", err)
		}
		defer svc.Close()

		topicID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Error converting topic ID to integer: %v", err)
		}

		fmt.Printf("Adding message to topic: %s with payload: %s\n", topicID, messagePayload)
		err = svc.PublishMessage(context.Background(), topicID, messagePayload, []byte{})
		if err != nil {
			log.Fatalf("Error adding message: %s", err)
		}
		fmt.Println("Message added successfully")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.AddCommand(addTopicCmd)
	addTopicCmd.Flags().StringVarP(&configFile, "config", "d", "", "Path to topic configuration file")

	addCmd.AddCommand(addSubscriptionCmd)
	addSubscriptionCmd.Flags().StringVarP(&configFile, "config", "d", "", "Path to subscription configuration file")

	addCmd.AddCommand(addMessageCmd)
	addMessageCmd.Flags().StringVarP(&messagePayload, "message", "d", "", "Message payload")
}
