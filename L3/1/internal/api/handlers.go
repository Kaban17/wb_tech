package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"wb_tech/l3_1/pkg/types"

	"github.com/gorilla/mux"
)

func (s *Server) handleGetNotify(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling get notify", "method", r.Method, "url", r.URL)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error("invalid notification id", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notification, err := s.repo.GetNotification(id)
	if err != nil {
		slog.Error("failed to get notification", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notification)
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

	if err := s.producer.Publish(notification, delay); err != nil {
		slog.Error("failed to publish notification", "error", err)
		http.Error(w, "failed to publish notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "notification accepted for processing"})
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
