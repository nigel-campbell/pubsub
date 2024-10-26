package pubsub

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	s, err := NewService(DefaultFilename)
	ok(t, err)
	defer s.Close()

	err = s.Init(ctx)
	ok(t, err)

	err = s.CreateTopic(ctx, "topic1", []byte("metadata"))
	ok(t, err)

	topic, err := s.GetTopic(ctx, "topic1")
	ok(t, err)
	equals(t, "topic1", topic.Name, "topic doesn't match expectation")

	err = s.CreateSubscription(ctx, topic.ID, "subscriber1", []byte("metadata"))
	ok(t, err)

	subscription, err := s.GetSubscription(ctx, topic.ID, "subscriber1")
	ok(t, err)
	equals(t, topic.ID, subscription.TopicID, "topic id doesn't match expectation")
	equals(t, "subscriber1", subscription.SubscriberID, "subscriber id doesn't match expectation")

	err = s.PublishMessage(ctx, topic.ID, "content", []byte("metadata"))
	ok(t, err)

	messages, err := s.GetMessages(ctx, subscription.ID)
	ok(t, err)
	equals(t, 1, len(messages), "message count doesn't match expectation")
	equals(t, "content", messages[0].Content, "message content doesn't match expectation")

	message := messages[0]

	now := time.Now()
	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err)
	equals(t, 1, len(messages), "message count when pulling messages doesn't match expectation")
	equals(t, "content", messages[0].Content, "message content when pulling messages doesn't match expectation")

	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err)
	equals(t, 0, len(messages), "message count when pulling messages again doesn't match expectation")

	err = s.ModifyAckDeadline(ctx, subscription.ID, message.ID, now.Add(-time.Minute*12))
	ok(t, err)

	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err)
	equals(t, 1, len(messages), "message count after modifying ack deadline doesn't match expectation")

	err = s.AcknowledgeMessage(ctx, subscription.ID, message.ID)
	ok(t, err)

	messages, err = s.PullMessages(ctx, subscription.ID, now.Add(time.Second*10))
	ok(t, err)
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

func ok(t *testing.T, err error) {
	if err != nil {
		_ = os.Remove("pubsub.db")
		t.Fatalf("error: %v", err)
	}
}
