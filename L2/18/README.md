# Calendar API Server

This is a simple HTTP server for managing a small calendar of events.

## Project Structure

- `cmd/calendar/main.go`: The main application entry point.
- `internal/app/api`: Contains the HTTP handlers and middleware.
- `internal/app/config`: Configuration loading.
- `internal/app/models`: Defines the `Event` data structure.
- `internal/app/storage`: In-memory storage for events and business logic.
- `go.mod`, `go.sum`: Go module files.

## API Endpoints

### Create Event

- **URL:** `/create_event`
- **Method:** `POST`
- **Body (JSON):**
  ```json
  {
    "user_id": 1,
    "date": "2023-12-31T00:00:00Z",
    "event": "New Year's Eve"
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "result": {
      "id": 1,
      "user_id": 1,
      "date": "2023-12-31T00:00:00Z",
      "event": "New Year's Eve"
    }
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: Invalid request body.
  - `503 Service Unavailable`: Date is already busy.
  - `500 Internal Server Error`: Other errors.

### Update Event

- **URL:** `/update_event`
- **Method:** `POST`
- **Body (JSON):**
  ```json
  {
    "id": 1,
    "user_id": 1,
    "date": "2023-12-31T00:00:00Z",
    "event": "Updated New Year's Eve"
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "result": {
      "id": 1,
      "user_id": 1,
      "date": "2023-12-31T00:00:00Z",
      "event": "Updated New Year's Eve"
    }
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: Invalid request body.
  - `503 Service Unavailable`: Event not found.
  - `500 Internal Server Error`: Other errors.

### Delete Event

- **URL:** `/delete_event`
- **Method:** `POST`
- **Body (JSON):**
  ```json
  {
    "id": 1
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "result": "event deleted"
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: Invalid request body.
  - `503 Service Unavailable`: Event not found.
  - `500 Internal Server Error`: Other errors.

### Get Events for Day

- **URL:** `/events_for_day`
- **Method:** `GET`
- **Query Parameters:**
  - `user_id`: (integer)
  - `date`: (string, `YYYY-MM-DD`)
- **Example:** `/events_for_day?user_id=1&date=2023-12-31`
- **Success Response (200 OK):**
  ```json
  {
    "result": [
      {
        "id": 1,
        "user_id": 1,
        "date": "2023-12-31T00:00:00Z",
        "event": "New Year's Eve"
      }
    ]
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: Invalid `user_id` or `date`.
  - `500 Internal Server Error`: Other errors.

### Get Events for Week

- **URL:** `/events_for_week`
- **Method:** `GET`
- **Query Parameters:**
  - `user_id`: (integer)
  - `date`: (string, `YYYY-MM-DD`)
- **Example:** `/events_for_week?user_id=1&date=2023-12-31`
- **Success Response (200 OK):**
  ```json
  {
    "result": [
      {
        "id": 1,
        "user_id": 1,
        "date": "2023-12-31T00:00:00Z",
        "event": "New Year's Eve"
      }
    ]
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: Invalid `user_id` or `date`.
  - `500 Internal Server Error`: Other errors.

### Get Events for Month

- **URL:** `/events_for_month`
- **Method:** `GET`
- **Query Parameters:**
  - `user_id`: (integer)
  - `date`: (string, `YYYY-MM-DD`)
- **Example:** `/events_for_month?user_id=1&date=2023-12-31`
- **Success Response (200 OK):**
  ```json
  {
    "result": [
      {
        "id": 1,
        "user_id": 1,
        "date": "2023-12-31T00:00:00Z",
        "event": "New Year's Eve"
      }
    ]
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: Invalid `user_id` or `date`.
  - `500 Internal Server Error`: Other errors.

## How to Run

1.  **Clone the repository.**
2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Run the server:**
    ```bash
    go run cmd/calendar/main.go
    ```
    The server will start on port `8080` by default.

4.  **Run tests:**
    ```bash
    go test ./...
    ```

## Configuration

The server can be configured using environment variables:

- `HTTP_PORT`: The port for the HTTP server to listen on (default: `8080`).
