package handlers

import (
	"net/http"
	"wb_tech/l3_3/pkg/service"

	"github.com/go-chi/chi/v5"
)

// API struct holds the service and router
type API struct {
	service *service.Service
	router  *chi.Mux
}

// NewAPI creates a new API instance with the given service
func NewAPI(service *service.Service) *API {
	return &API{
		service: service,
		router:  chi.NewRouter(),
	}
}

// Router sets up all the API routes and returns the http.Handler
func (a *API) Router() http.Handler {
	a.router.Post("/comments", a.CreateComment)
	a.router.Get("/comments", a.GetComments) // Changed from /comments/search
	a.router.Delete("/comments/{id}", a.DeleteComment)
	a.router.Get("/ping", a.Ping)
	a.router.Get("/", a.GetFrontEnd)

	// Serve static files from the /web directory
	fileServer := http.FileServer(http.Dir("/web"))
	a.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return a.router
}

// GetFrontEnd serves the index.html file
func (a *API) GetFrontEnd(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/web/index.html")
}
