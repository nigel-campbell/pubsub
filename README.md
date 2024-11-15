# PubSub CLI

A command-line tool to emulate a Pub/Sub service. This tool allows you to manage topics, subscriptions, and messages, mimicking the behavior of a Pub/Sub system in a local SQLite database.

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

## License

This project is licensed under the terms of the [MIT License](LICENSE).
