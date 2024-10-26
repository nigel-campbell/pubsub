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

// ackCmd represents the "ack" command
var ackCmd = &cobra.Command{
	Use:   "ack [SUBSCRIPTION_ID] [MESSAGE_ID]",
	Short: "Acknowledge a message in a subscription",
	Long:  "Marks a message as acknowledged in a specific subscription.",
	Args:  cobra.ExactArgs(2),
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

		messageID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Invalid message ID: %s", args[1])
		}

		err = svc.AcknowledgeMessage(ctx, subscriptionID, messageID)
		if err != nil {
			log.Fatalf("Failed to acknowledge message %d in subscription %d: %v", messageID, subscriptionID, err)
		}
		fmt.Printf("Acknowledged message %d in subscription %d\n", messageID, subscriptionID)
	},
}

var modAckCmd = &cobra.Command{
	Use:   "modack [SUBSCRIPTION_ID] [MESSAGE_ID] [DEADLINE]",
	Short: "Modify the ack deadline for a message",
	Long:  "Modifies the ack deadline for a message in a specific subscription",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Error creating Pub/Sub service: %v", err)
		}
		defer svc.Close()

		subscriptionId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid subscription ID: %s", args[0])
		}

		messageID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Invalid message ID: %s", args[1])
		}

		deadline, err := time.ParseDuration(args[2])
		if err != nil {
			log.Fatalf("Invalid deadline: %s", args[2])
		}

		err = svc.ModifyAckDeadline(ctx, subscriptionId, messageID, time.Now().Add(deadline))
		if err != nil {
			log.Fatalf("Failed to modify ack deadline for message %d: %v", messageID, err)
			log.Fatalf("Failed to modify ack deadline for message %d: %v", messageID, err)
			log.Fatalf("Failed to modify ack deadline for message %d: %v", messageID, err)
		}
		fmt.Printf("Modified ack deadline for message %d\n", messageID)
	},
}

var nackCmd = &cobra.Command{
	Use:   "nack [SUBSCRIPTION_ID] [MESSAGE_ID]",
	Short: "Modify the ack deadline for a message",
	Long:  "Modifies the ack deadline for a message in a specific subscription",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		svc, err := pubsub.NewService(pubsub.DefaultFilename)
		if err != nil {
			log.Fatalf("Error creating Pub/Sub service: %v", err)
		}
		defer svc.Close()

		subscriptionId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid subscription ID: %s", args[0])
		}

		messageID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Invalid message ID: %s", args[1])
		}

		// NB: This could be implemented in modack but cobra doesn't play well with negative numbers since the negative
		// sign is interpreted as a flag. This is a workaround.
		err = svc.ModifyAckDeadline(ctx, subscriptionId, messageID, time.Now())
		if err != nil {
			log.Fatalf("Failed to modify ack deadline for message %d: %v", messageID, err)
		}
		fmt.Printf("Modified ack deadline for message %d\n", messageID)
	},
}

func init() {
	rootCmd.AddCommand(ackCmd)
	rootCmd.AddCommand(modAckCmd)
	rootCmd.AddCommand(nackCmd)
}
