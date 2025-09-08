#!/bin/bash

# This script starts the server, runs a series of API tests using curl,
# and then stops the server.

# Exit immediately if a command exits with a non-zero status.
set -e

echo "Starting server in background..."
go run cmd/calendar/main.go &
# Get the process ID of the server
SERVER_PID=$!

# Give the server a moment to start
sleep 2

echo "Server started with PID $SERVER_PID. Running tests..."

# Function to check if the last command was successful
check_success() {
  if [ $? -eq 0 ]; then
    echo "SUCCESS"
  else
    echo "FAILURE"
    # Clean up and exit if any command fails
    kill $SERVER_PID
    exit 1
  fi
}

# Test 1: Create Event
echo -n "Test 1: POST /create_event... "
CREATE_RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d '{"user_id": 10, "date": "2025-10-01T12:00:00Z", "event": "API Test Event"}' http://localhost:8080/create_event)
echo "$CREATE_RESPONSE" | grep -q '"id":'
check_success

# Extract event ID from the response
EVENT_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)

# Test 2: Get Events for Day
echo -n "Test 2: GET /events_for_day... "
curl -s "http://localhost:8080/events_for_day?user_id=10&date=2025-10-01" | grep -q "API Test Event"
check_success

# Test 3: Get Events for Week
echo -n "Test 3: GET /events_for_week... "
curl -s "http://localhost:8080/events_for_week?user_id=10&date=2025-10-01" | grep -q "API Test Event"
check_success

# Test 4: Get Events for Month
echo -n "Test 4: GET /events_for_month... "
curl -s "http://localhost:8080/events_for_month?user_id=10&date=2025-10-01" | grep -q "API Test Event"
check_success

# Test 5: Update Event
echo -n "Test 5: POST /update_event... "
curl -s -X POST -H "Content-Type: application/json" -d "{\"id\": $EVENT_ID, \"user_id\": 10, \"date\": \"2025-10-01T12:00:00Z\", \"event\": \"Updated API Test Event\"}" http://localhost:8080/update_event | grep -q "Updated API Test Event"
check_success

# Test 6: Delete Event
echo -n "Test 6: POST /delete_event... "
curl -s -X POST -H "Content-Type: application/json" -d "{\"id\": $EVENT_ID}" http://localhost:8080/delete_event | grep -q "event deleted"
check_success

echo "All API tests passed!"

echo "Stopping server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

echo "Server stopped."
