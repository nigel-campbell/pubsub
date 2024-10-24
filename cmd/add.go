/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Define variables to store flags (like -d for config or payload)
var configFile string
var messagePayload string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add topics, subscriptions, or messages",
	Long:  "Use this command to add new topics, subscriptions, or messages to the Pub/Sub system",
}

// addTopicCmd represents the "add topic" command
var topicCmd = &cobra.Command{
	Use:   "topic [TOPIC_ID]",
	Short: "Add a new topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topicID := args[0]
		// Implement the logic for adding a topic using the topicID and configFile
		fmt.Printf("Adding topic: %s with config: %s\n", topicID, configFile)
		// Call your internal logic here to add the topic
	},
}

// addSubscriptionCmd represents the "add subscription" command
var subscriptionCmd = &cobra.Command{
	Use:   "subscription [TOPIC_ID] [SUBSCRIPTION_ID]",
	Short: "Add a new subscription to a topic",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		topicID := args[0]
		subscriptionID := args[1]
		// Implement the logic for adding a subscription using the topicID, subscriptionID, and configFile
		fmt.Printf("Adding subscription: %s to topic: %s with config: %s\n", subscriptionID, topicID, configFile)
		// Call your internal logic here to add the subscription
	},
}

// addMessageCmd represents the "add message" command
var messageCmd = &cobra.Command{
	Use:   "message [TOPIC_ID]",
	Short: "Add a message to a topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topicID := args[0]
		// Implement the logic for adding a message using the topicID and messagePayload
		fmt.Printf("Adding message to topic: %s with payload: %s\n", topicID, messagePayload)
		// Call your internal logic here to add the message
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// Register "add topic" command
	addCmd.AddCommand(topicCmd)
	topicCmd.Flags().StringVarP(&configFile, "config", "d", "", "Path to topic configuration file")
	topicCmd.MarkFlagRequired("config") // Make sure config flag is required

	// Register "add subscription" command
	addCmd.AddCommand(subscriptionCmd)
	subscriptionCmd.Flags().StringVarP(&configFile, "config", "d", "", "Path to subscription configuration file")
	subscriptionCmd.MarkFlagRequired("config") // Make sure config flag is required

	// Register "add message" command
	addCmd.AddCommand(messageCmd)
	messageCmd.Flags().StringVarP(&messagePayload, "message", "d", "", "Message payload")
	messageCmd.MarkFlagRequired("message") // Make sure message payload flag is required
}
