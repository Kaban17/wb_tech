package handlers

import (
	"net/http"
	"wb_tech/l3_3/pkg/service"

	"github.com/go-chi/chi/v5"
)

type API struct {
	service *service.Service
	router  *chi.Mux
}

func (a *API) GetFrontEnd(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/web/index.html")
}

func NewAPI(service *service.Service) *API {
	return &API{
		service: service,
		router:  chi.NewRouter(),
	}
}
func (a *API) Router() http.Handler {
	a.router.Post("/comments", a.CreateComment)
	a.router.Get("/comments/search", a.GetComments)
	a.router.Get("/ping", a.Ping)
	a.router.Get("/", a.GetFrontEnd)

	fileServer := http.FileServer(http.Dir("/web"))
	a.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return a.router
}
