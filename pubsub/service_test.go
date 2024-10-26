package pubsub

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	_ = os.Remove(DefaultFilename)

	s, err := NewService(DefaultFilename)
	ok(t, err, "failed to create service")
	defer s.Close()

	err = s.Init(ctx)
	ok(t, err, "failed to initialize service")

	err = s.CreateTopic(ctx, "topic1", []byte("metadata"))
	ok(t, err, "failed to create topic")

	topic, err := s.GetTopic(ctx, "topic1")
	ok(t, err, "failed to get topic")
	equals(t, "topic1", topic.Name, "topic doesn't match expectation")

	err = s.CreateSubscription(ctx, topic.ID, "subscriber1", []byte("metadata"))
	ok(t, err, "failed to create subscription")

	subscription, err := s.GetSubscription(ctx, topic.ID, "subscriber1")
	ok(t, err, "failed to get subscription")
	equals(t, topic.ID, subscription.TopicID, "topic id doesn't match expectation")
	equals(t, "subscriber1", subscription.SubscriberID, "subscriber id doesn't match expectation")

	err = s.PublishMessage(ctx, topic.ID, "content", []byte("metadata"))
	ok(t, err, "failed to publish message")

	messages, err := s.GetMessages(ctx, subscription.ID)
	ok(t, err, "failed to get messages")
	equals(t, 1, len(messages), "message count doesn't match expectation")
	equals(t, "content", messages[0].Content, "message content doesn't match expectation")

	message := messages[0]

	now := time.Now()
	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err, "failed to pull messages")
	equals(t, 1, len(messages), "message count when pulling messages doesn't match expectation")
	equals(t, "content", messages[0].Content, "message content when pulling messages doesn't match expectation")

	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err, "failed to pull messages again")
	equals(t, 0, len(messages), "message count when pulling messages again doesn't match expectation")

	err = s.ModifyAckDeadline(ctx, subscription.ID, message.ID, now.Add(-time.Minute*12))
	ok(t, err, "failed to modify ack deadline")

	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err, "failed to pull messages after modifying ack deadline")
	equals(t, 1, len(messages), "message count after modifying ack deadline doesn't match expectation")

	err = s.AcknowledgeMessage(ctx, subscription.ID, message.ID)
	ok(t, err, "failed to acknowledge message")

	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err, "failed to pull messages after acknowledging message")
	equals(t, 0, len(messages), "message count after acknowledging message doesn't match expectation")

	_ = os.Remove(DefaultFilename)
	t.Log("Test successful. Database removed")
}

func equals(t *testing.T, expected, actual interface{}, desc string) {
	if expected != actual {
		_ = os.Remove("pubsub.db")
		t.Fatalf("%s: expected %v, got %v", desc, expected, actual)
	}
}

func ok(t *testing.T, err error, desc string) {
	if err != nil {
		_ = os.Remove("pubsub.db")
		t.Fatalf("error: %v", err)
	}
}
