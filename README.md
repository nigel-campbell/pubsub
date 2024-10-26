# PubSub CLI

A command-line tool to emulate a Pub/Sub service. This tool allows you to manage topics, subscriptions, and messages, mimicking the behavior of a Pub/Sub system in a local SQLite database.

Intended for https://github.com/nigel-campbell/pubsub-emulator.

## Features

- **Initialize** the Pub/Sub environment with necessary tables and database setup.
- **Add** topics, subscriptions, and messages.
- **List** topics, subscriptions, and messages.
- **Acknowledge (Ack)** messages to mark them as processed.
- **Pull** unacknowledged messages for consumption.
- **Clean** the database to remove all topics, subscriptions, and messages.

## Installation

1. Clone the repository:
   ```bash
   git clone github.com/nigel-campbell/pubsub.git
   cd pubsub
   ```

2. Build the CLI:
   ```bash
   make build
   ```

3. Run the CLI:
   ```bash
   ./bin/pubsub
   ```

## Usage

The basic commands are:

```bash
./bin/pubsub init                                  # Initialize database and tables
./bin/pubsub add topic <TOPIC_NAME> -d <CONFIG>    # Add a topic
./bin/pubsub add subscription <TOPIC_ID> <SUBSCRIPTION_ID> -d <CONFIG>   # Add a subscription
./bin/pubsub add message <TOPIC_ID> -d <MESSAGE_PAYLOAD>                # Add a message
./bin/pubsub list topics                           # List all topics
./bin/pubsub list subscriptions <TOPIC_ID>         # List subscriptions for a topic
./bin/pubsub pull <SUBSCRIPTION_ID>                # Pull messages for a subscription
./bin/pubsub ack <SUBSCRIPTION_ID> <MESSAGE_ID>    # Acknowledge a message
./bin/pubsub clean                                 # Clean all data
```

## Project Structure

```plaintext
├── LICENSE             # License file
├── Makefile            # Build and setup commands
├── README.md           # Project documentation
├── bin/                # Compiled binary
│   └── pubsub
├── cmd/                # CLI command implementations
│   ├── ack.go
│   ├── add.go
│   ├── clean.go
│   ├── init.go
│   ├── list.go
│   ├── pull.go
│   ├── root.go
│   ├── subscriptions.go
│   └── topics.go
├── go.mod              # Go module file
├── go.sum              # Go dependency file
├── pubsub/             # Service and database logic
│   ├── service.go
│   └── service_test.go
├── pubsub.db           # SQLite database file
└── pubsub.go           # Main application entry point
```

## License

This project is licensed under the terms of the [MIT License](LICENSE).