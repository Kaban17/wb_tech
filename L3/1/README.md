# Notification Service

This is a simple notification service written in Go. It allows scheduling notifications to be sent at a later time. The service uses RabbitMQ for message queuing with delayed delivery and PostgreSQL for storing notification data. It exposes a RESTful API for creating and checking the status of notifications, using Server-Sent Events (SSE) to provide real-time updates.

## Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.20+ recommended)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repo-url>
    cd <project-directory>
    ```

2.  **Set up environment variables:**
    Create a `.env` file in the root of the project by copying the example:
    ```bash
    cp .env.example .env
    ```
    Update the `.env` file with your configuration for PostgreSQL and RabbitMQ. A typical configuration might look like this:
    ```
    # PostgreSQL settings
    DB_USER=user
    DB_PASSWORD=password
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=notifications
    DB_SSLMODE=disable

    # RabbitMQ settings
    RABBITMQ_URL=amqp://guest:guest@localhost:5672/
    ```

3.  **Start dependent services:**
    A `docker-compose.yml` file is included to easily start PostgreSQL and RabbitMQ.
    ```bash
    docker-compose up -d
    ```

## Running the Application

1.  **Build and run the service:**
    The included `Makefile` simplifies the process. To build and run the application, simply use:
    ```bash
    make run
    ```
    The server will start on port `8080`.

2.  **Static Analysis:**
    To run linters and vet checks, use:
    ```bash
    make all
    ```

3.  **Run Tests:**
    To run the test suite:
    ```bash
    make test
    ```

## API Documentation

### Create a Notification

-   **Endpoint:** `POST /notify`
-   **Method:** `POST`
-   **Description:** Creates and schedules a new notification.
-   **Query Parameters:**
    -   `delay` (optional): The delay in milliseconds before the notification is sent. If omitted, the notification is sent immediately.
-   **Request Body:** A JSON object representing the notification.
    ```json
    {
      "message": "Your message content",
      "mail": "recipient@example.com",
      "tg": "@telegram_user"
    }
    ```
-   **Example Request (`curl`):**
    ```bash
    # Schedule a notification with a 10-second delay
    curl -X POST 'http://localhost:8080/notify?delay=10000' \
    -H "Content-Type: application/json" \
    -d '{"message": "Hello, this is a delayed notification!", "mail": "test@example.com", "tg": "@test_user"}'
    ```
-   **Success Response:**
    -   **Code:** `202 Accepted`
    -   **Content:**
        ```json
        {
          "id": "1",
          "message": "notification accepted for processing with id"
        }
        ```

### Get Notification Status (SSE)

-   **Endpoint:** `GET /notify/{id}`
-   **Method:** `GET`
-   **Description:** Retrieves the status of a notification using a Server-Sent Events (SSE) connection. The connection remains open and pushes updates until the notification is delivered.
-   **URL Parameters:**
    -   `id` (required): The ID of the notification to track.
-   **Example Request (`curl`):**
    ```bash
    curl -H "Accept: text/event-stream" http://localhost:8080/notify/1
    ```
-   **Response Stream:**
    The server will send events as the notification status changes.
    ```
    event: initial
    data: {"id":1,"status":"pending",...}

    event: update
    data: {"id":1,"status":"delivered",...}
    ```
    If the notification has already been delivered when the request is made, the server will send the final status and immediately close the connection.
