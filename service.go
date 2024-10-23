package main

import (
	"context"
	"database/sql"
)

const (
	// SQL statements to create tables
	createTopicsTable = `CREATE TABLE IF NOT EXISTS Topics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        metadata BLOB,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	createSubscriptionsTable = `CREATE TABLE IF NOT EXISTS Subscriptions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        topic_id INTEGER NOT NULL,
        subscriber_id TEXT NOT NULL,
        metadata BLOB,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (topic_id) REFERENCES Topics(id)
    );`

	createMessagesTable = `CREATE TABLE IF NOT EXISTS Messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        topic_id INTEGER NOT NULL,
        subscription_id INTEGER,
        content TEXT NOT NULL,
        metadata BLOB,
        published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (topic_id) REFERENCES Topics(id),
        FOREIGN KEY (subscription_id) REFERENCES Subscriptions(id)
    );`
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
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTopic(ctx context.Context, name string, metadata []byte) error {
	_, err := s.db.Exec("INSERT INTO Topics (name, metadata) VALUES (?, ?)", name, metadata)
	return err
}

func (s *Service) GetTopic(ctx context.Context, name string) (*Topic, error) {
	row := s.db.QueryRow("SELECT id, name, metadata FROM Topics WHERE name = ?", name)
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

func (s *Service) PublishMessage(ctx context.Context, topicID int, subscriptionID int, content string, metadata []byte) error {
	_, err := s.db.Exec("INSERT INTO Messages (topic_id, subscription_id, content, metadata) VALUES (?, ?, ?, ?)", topicID, subscriptionID, content, metadata)
	return err
}

func (s *Service) GetMessages(ctx context.Context, topicID int, subscriptionID int) ([]*Message, error) {
	rows, err := s.db.Query("SELECT id, topic_id, subscription_id, content, metadata FROM Messages WHERE topic_id = ? AND subscription_id = ?", topicID, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}
	for rows.Next() {
		message := &Message{}
		if err := rows.Scan(&message.ID, &message.TopicID, &message.SubscriptionID, &message.Content, &message.Metadata); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (s *Service) Init() error {
	_, err := s.db.Exec(createTopicsTable)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(createSubscriptionsTable)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(createMessagesTable)
	if err != nil {
		return err
	}
	return nil
}
