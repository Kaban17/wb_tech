package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wb_tech/l3_3/pkg/repository"
)

func (a *API) CreateComment(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateComment"

	var comment repository.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := a.service.Add(comment, comment.ParentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (a *API) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
func (a *API) GetComments(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetComments"
	query := r.URL.Query().Get("query")
	parent := r.URL.Query().Get("parent")

	if query != "" {
		comments, err := a.service.Search(query, 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	var parentID *int
	if parent != "" {
		id, err := strconv.Atoi(parent)
		if err != nil {
			http.Error(w, "invalid parent id", http.StatusBadRequest)
			return
		}
		parentID = &id
	}

	// TODO: get limit, offset, sortField, sortOrder from query params
	comments, err := a.service.GetThread(parentID, 100, 0, "path", "asc")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
