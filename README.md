# Go Message Broker

A simple message broker implemented in Go.

## Features

- In-memory message queues
- Fixed-size queues
- Limited number of subscribers per queue
- JSON-based messages
- Simple HTTP API

## API

### Send a message

- **URL:** `/v1/queues/{queue_name}/messages`
- **Method:** `POST`
- **Body:** JSON message
- **Example:**
  ```bash
  curl -X POST -d '{"event":"delivered"}' http://localhost:8080/v1/queues/app_events/messages
  ```

### Subscribe to a queue

- **URL:** `/v1/queues/{queue_name}/subscriptions`
- **Method:** `POST`
- **Description:** Subscribes to a queue and receives messages as they are sent. The connection is kept open.
- **Example:**
  ```bash
  curl -X POST http://localhost:8080/v1/queues/app_events/subscriptions
  ```

## How to run

### Locally

1.  Start the broker:
    ```bash
    go run ./cmd/broker/main.go
    ```
2.  In a new terminal, subscribe to a queue:
    ```bash
    go run ./subscriber/main.go
    ```
3.  In another terminal, send a message:
    ```bash
    curl -X POST -d '{"event":"delivered"}' http://localhost:8080/v1/queues/app_events/messages
    ```

### With Docker

1.  Build the Docker image:
    ```bash
    docker build -t broker .
    ```
2.  Run the Docker container:
    ```bash
    docker run -p 8080:8080 broker
    ```

