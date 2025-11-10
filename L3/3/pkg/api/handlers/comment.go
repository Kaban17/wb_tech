package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wb_tech/l3_3/pkg/repository"

	"github.com/go-chi/chi/v5"
)

// The API struct is defined in router.go. We do not redefine it here to avoid redeclaration errors.

// CreateComment handles the creation of a new comment
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

// Ping handles the /ping endpoint
func (a *API) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// GetComments handles fetching comments based on query or parent ID
func (a *API) GetComments(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetComments"
	query := r.URL.Query().Get("query")
	parent := r.URL.Query().Get("parent")

	var comments []repository.Comment
	var err error

	if query != "" {
		comments, err = a.service.Search(query, 10, 0)
	} else {
		var parentID *int
		if parent != "" {
			id, parseErr := strconv.Atoi(parent)
			if parseErr != nil {
				http.Error(w, "invalid parent id", http.StatusBadRequest)
				return
			}
			parentID = &id
		}

		// TODO: get limit, offset, sortField, sortOrder from query params
		comments, err = a.service.GetThread(parentID, 100, 0, "path", "asc")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure an empty slice is returned as [] not null
	if comments == nil {
		comments = []repository.Comment{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteComment handles the deletion of a comment and its subtree
func (a *API) DeleteComment(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DeleteComment"

	commentIDStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := a.service.DeleteSubtree(commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}
