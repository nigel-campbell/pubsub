package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	db *sql.DB
}

type Topic struct {
	ID       int
	Name     string
	Metadata []byte
}

type Subscription struct {
	ID           int
	TopicID      int
	SubscriberID string
}

type Message struct {
	ID             int
	TopicID        int
	SubscriptionID int
	Content        string
	Metadata       []byte
	Acknowledged   bool
	AckDeadline    sql.NullTime // Use sql.NullTime for fields that may not always have a value
}

func (m *Message) String() string {
	return fmt.Sprintf("ID: %d, TopicID: %d, SubscriptionID: %d, Content: %s, Metadata: %s, Acknowledged: %t, AckDeadline: %v", m.ID, m.TopicID, m.SubscriptionID, m.Content, string(m.Metadata), m.Acknowledged, m.AckDeadline)
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTopic(ctx context.Context, name string, metadata []byte) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO Topics (name, metadata) VALUES (?, ?)", name, metadata)
	return err
}

func (s *Service) GetTopic(ctx context.Context, name string) (*Topic, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, name, metadata FROM Topics WHERE name = ?", name)
	topic := &Topic{}
	err := row.Scan(&topic.ID, &topic.Name, &topic.Metadata)
	return topic, err
}

func (s *Service) CreateSubscription(ctx context.Context, topicID int, subscriberID string, metadata []byte) error {
	_, err := s.db.Exec("INSERT INTO Subscriptions (topic_id, subscriber_id, metadata) VALUES (?, ?, ?)", topicID, subscriberID, metadata)
	return err
}

func (s *Service) GetSubscription(ctx context.Context, topicID int, subscriberID string) (*Subscription, error) {
	row := s.db.QueryRow("SELECT id, topic_id, subscriber_id FROM Subscriptions WHERE topic_id = ? AND subscriber_id = ?", topicID, subscriberID)
	subscription := &Subscription{}
	err := row.Scan(&subscription.ID, &subscription.TopicID, &subscription.SubscriberID)
	return subscription, err
}

func (s *Service) PublishMessage(ctx context.Context, topicID int, content string, metadata []byte) error {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rows, err := tx.QueryContext(ctx, "SELECT id FROM Subscriptions WHERE topic_id = ?", topicID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("query error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var subscriptionID int
		if err := rows.Scan(&subscriptionID); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("scan error: %v, rollback error: %v", err, rbErr)
			}
			return err
		}

		_, err := tx.ExecContext(ctx, "INSERT INTO Messages (topic_id, subscription_id, content, metadata) VALUES (?, ?, ?, ?)",
			topicID, subscriptionID, content, metadata)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("insert error: %v, rollback error: %v", err, rbErr)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetMessages returns all messages for a subscription regardless of acknowledgement status
func (s *Service) GetMessages(ctx context.Context, subscriptionID int) ([]*Message, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, topic_id, subscription_id, content, metadata, acknowledged, ack_deadline FROM Messages WHERE subscription_id = ?", subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		if err := rows.Scan(&message.ID, &message.TopicID, &message.SubscriptionID, &message.Content, &message.Metadata, &message.Acknowledged, &message.AckDeadline); err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate messages: %v", err)
	}

	return messages, nil
}

// PullMessages returns all messages that have not been acknowledged and have not passed their ack_deadline.
func (s *Service) PullMessages(ctx context.Context, subscriptionID int, ackDeadline time.Time) ([]*Message, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT id, topic_id, subscription_id, content, metadata, acknowledged, ack_deadline FROM Messages WHERE subscription_id = ? AND acknowledged = 0", subscriptionID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("query error: %v, rollback error: %v", err, rbErr)
		}
		return nil, fmt.Errorf("failed to pull messages: %v", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		if err := rows.Scan(&message.ID, &message.TopicID, &message.SubscriptionID, &message.Content, &message.Metadata, &message.Acknowledged, &message.AckDeadline); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return nil, fmt.Errorf("scan error: %v, rollback error: %v", err, rbErr)
			}
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}

		now := time.Now()
		if now.Before(message.AckDeadline.Time) {
			// To reduce message redelivery.
			continue
		}
		messages = append(messages, message)
	}

	for _, message := range messages {
		_, err := tx.ExecContext(ctx, "UPDATE Messages SET ack_deadline = ? WHERE id = ?", ackDeadline, message.ID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return nil, fmt.Errorf("update error: %v, rollback error: %v", err, rbErr)
			}
			return nil, fmt.Errorf("failed to update message: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("rows error: %v, rollback error: %v", err, rbErr)
		}
		return nil, fmt.Errorf("failed to iterate messages: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction when pulling messages: %v", err)
	}

	return messages, nil
}

// AcknowledgeMessage sets the acknowledged field to true for a message
func (s *Service) AcknowledgeMessage(ctx context.Context, messageID int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE Messages SET acknowledged = 1 WHERE id = ?", messageID)
	return err
}

func (s *Service) ModifyAckDeadline(ctx context.Context, messageID int, ackDeadline time.Time) error {
	_, err := s.db.ExecContext(ctx, "UPDATE Messages SET ack_deadline = ? WHERE id = ?", ackDeadline, messageID)
	return err
}

func (s *Service) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Topics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        metadata BLOB,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Subscriptions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        topic_id INTEGER NOT NULL,
        subscriber_id TEXT NOT NULL,
        metadata BLOB,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (topic_id) REFERENCES Topics(id)
    );`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        topic_id INTEGER NOT NULL,
        subscription_id INTEGER,
        content TEXT NOT NULL,
        metadata BLOB,
        acknowledged BOOLEAN DEFAULT FALSE,
        ack_deadline DATETIME,
        published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (topic_id) REFERENCES Topics(id),
        FOREIGN KEY (subscription_id) REFERENCES Subscriptions(id)
    );`)
	if err != nil {
		return err
	}
	return nil
}
