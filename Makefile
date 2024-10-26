# Variables
BINARY_NAME := pubsub
BUILD_DIR := ./bin
ENTRY_POINT := ./pubsub.go

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(ENTRY_POINT)

# Clean up binary and other generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Help message
help:
	@echo "Usage:"
	@echo "  make build      Build the binary"
	@echo "  make run        Run the binary"
	@echo "  make clean      Remove built files"

.PHONY: all build run clean help
