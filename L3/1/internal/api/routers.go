package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/notify/{id}", handleGetNotify).Methods("GET")
	r.HandleFunc("/notify", handlePostNotify).Methods("POST")
	r.HandleFunc("/notify/{id}", handlePutNotify).Methods("PUT")
	r.HandleFunc("/notify/{id}", handleDeleteNotify).Methods("DELETE")

	return r
}
