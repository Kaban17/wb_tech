package api

import (
	"database/sql"
	"net/http"
	"wb_tech/l3_1/internal/storage/postgres"
	"wb_tech/l3_1/internal/storage/queue"

	"github.com/gorilla/mux"
)

type Server struct {
	repo     *postgres.Repository
	producer *queue.Producer
}

func NewServer(db *sql.DB, producer *queue.Producer) *Server {
	return &Server{
		repo:     postgres.NewRepository(db),
		producer: producer,
	}
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
