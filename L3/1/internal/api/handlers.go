package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"wb_tech/l3_1/pkg/types"

	"github.com/gorilla/mux"
)

func (s *Server) handleGetNotify(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling get notify SSE", "method", r.Method, "url", r.URL)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("invalid notification id", "error", err)
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := s.repo.GetNotification(id)
	if err != nil {
		slog.Error("notification not found", "error", err)
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	// If the notification has already been delivered, send it and close the connection.
	if notification.Status == types.Delivered {
		s.sendSSEEvent(w, "final", notification)
		slog.Info("sent final status and closed connection", "notification_id", id)
		return
	}

	notifyChan := make(chan types.Notification, 10)
	s.registerNotification(id, notifyChan)
	defer s.unregisterNotification(id, notifyChan)

	s.sendSSEEvent(w, "initial", notification)

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	ctx := r.Context()

	for {
		select {
		case notification := <-notifyChan:
			s.sendSSEEvent(w, "update", notification)
			// If this is the final status, we can close the connection.
			if notification.Status == types.Delivered {
				slog.Info("sent final status update and closing connection", "notification_id", id)
				return
			}

		case <-heartbeat.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

		case <-ctx.Done():
			slog.Info("SSE connection closed by client", "notification_id", id)
			return
		}
	}
}

func (s *Server) sendSSEEvent(w http.ResponseWriter, eventType string, data interface{}) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        slog.Error("failed to marshal SSE data", "error", err)
        return
    }

    // Форматируем сообщение по спецификации SSE
    if eventType != "" {
        fmt.Fprintf(w, "event: %s\n", eventType)
    }
    fmt.Fprintf(w, "data: %s\n\n", jsonData)

    slog.Info("sending SSE event", "event", eventType, "data", string(jsonData))

    if flusher, ok := w.(http.Flusher); ok {
        flusher.Flush()
    }
}
func (s *Server) handlePostNotify(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling post notify", "method", r.Method, "url", r.URL)
	var notification *types.Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		slog.Error("failed to decode notification", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	delayStr := r.URL.Query().Get("delay")
	var delay time.Duration
	if delayStr != "" {
		parsedDelay, err := time.ParseDuration(delayStr + "ms")
		if err != nil {
			slog.Warn("invalid delay value, using 0", "delay", delayStr, "error", err)
			delay = 0
		} else {
			delay = parsedDelay
		}
	}
	now := time.Now()
	notification.TimeCreated = now
	notification.ScheduledAt = now.Add(delay)
	notification.TimeSent = now.Add(delay)
	notification.Status = types.Pending

	// 1. Create notification in the database to get an ID
	id, err := s.repo.CreateNotification(notification)
	if err != nil {
		slog.Error("failed to create notification", "error", err)
		http.Error(w, "failed to create notification", http.StatusInternalServerError)
		return
	}
	notification.ID = id

	// 2. Publish the notification with the ID to the queue
	if err := s.producer.Publish(notification); err != nil {
		slog.Error("failed to publish notification", "error", err)
		// Optionally, you might want to delete the created notification or mark it as failed
		http.Error(w, "failed to publish notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "notification accepted for processing with id", "id": strconv.Itoa(id)})
}

func (s *Server) handlePutNotify(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling put notify", "method", r.Method, "url", r.URL)
	var notification *types.Notification
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("invalid notification id", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		slog.Error("failed to decode notification", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.repo.UpdateNotification(notification, id); err != nil {
		slog.Error("failed to update notification", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteNotify(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling delete notify", "method", r.Method, "url", r.URL)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("invalid notification id", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.repo.DeleteNotification(id); err != nil {
		slog.Error("failed to delete notification", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
