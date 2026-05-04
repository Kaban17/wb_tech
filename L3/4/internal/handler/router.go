package handler

import (
	"net/http"
)

type ImageHandler interface {
	UploadImage(w http.ResponseWriter, r *http.Request)
	GetImageByID(w http.ResponseWriter, r *http.Request)
	DeleteImage(w http.ResponseWriter, r *http.Request)
}

func NewRouter(mux *http.ServeMux, imageHandler ImageHandler, staticDir string) {
	mux.HandleFunc("POST /upload", imageHandler.UploadImage)
	mux.HandleFunc("GET /image/{id}", imageHandler.GetImageByID)
	mux.HandleFunc("DELETE /image/{id}", imageHandler.DeleteImage)

	// Serve static files
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("GET /", fileServer)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Serve processed images
	mux.Handle("GET /data/", http.StripPrefix("/data/", http.FileServer(http.Dir("./data"))))
}
