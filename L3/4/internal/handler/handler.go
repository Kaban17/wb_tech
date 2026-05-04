package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"wb_tech/l3_4/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	SaveImage(ctx context.Context, id string, file io.Reader) error
	GetImageByID(ctx context.Context, id string) (*model.Image, error)
	DeleteImage(ctx context.Context, id string) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	id := uuid.New().String()
	if err := h.service.SaveImage(r.Context(), id, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *handler) GetImageByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	img, err := h.service.GetImageByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(img)
}

func (h *handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteImage(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
