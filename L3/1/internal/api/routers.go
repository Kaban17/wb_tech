package api

import (
	"database/sql"
	"log/slog"
	"net/http"
	"sync"
	"wb_tech/l3_1/internal/storage/postgres"
	"wb_tech/l3_1/internal/storage/queue"
	"wb_tech/l3_1/pkg/types"

	"github.com/gorilla/mux"
)

type Server struct {
	repo                 *postgres.Repository
	producer             *queue.Producer
	mu                   sync.RWMutex
	notificationChannels map[int]chan types.Notification // Добавляем карту каналов
}

func NewServer(db *sql.DB, producer *queue.Producer) *Server {
	return &Server{
		repo:                 postgres.NewRepository(db),
		producer:             producer,
		mu:                   sync.RWMutex{},
		notificationChannels: make(map[int]chan types.Notification),
	}
}

func (s *Server) SetProducer(p *queue.Producer) {
	s.producer = p
}

func NewRouter(s *Server) *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix("/notify").Subrouter()
	r.HandleFunc("/{id}", s.handleGetNotify).Methods("GET")
	r.HandleFunc("", s.handlePostNotify).Methods("POST")
	r.HandleFunc("/{id}", s.handlePutNotify).Methods("PUT")
	r.HandleFunc("/{id}", s.handleDeleteNotify).Methods("DELETE")

	return r
}

func Run(r *mux.Router, port string) error {
	return http.ListenAndServe(port, r)
}
func (s *Server) registerNotification(id int, ch chan types.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notificationChannels[id]; !ok {
		s.notificationChannels[id] = ch
	}
}

func (s *Server) unregisterNotification(id int, ch chan types.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notificationChannels[id]; ok {
		close(s.notificationChannels[id])
		delete(s.notificationChannels, id)
	}
}
func (s *Server) BroadcastNotification(notification types.Notification, id int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if ch, exists := s.notificationChannels[id]; exists {
		select {
		case ch <- notification:
			slog.Info("notification sent via SSE", "id", id)
		default:
			slog.Warn("notification channel blocked", "id", id)
		}
	}
}
