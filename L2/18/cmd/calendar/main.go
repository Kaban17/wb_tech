package main

import (
	"fmt"
	"log"
	"net/http"

	"calendar/internal/app/api"
	"calendar/internal/app/config"
	"calendar/internal/app/storage"
)

func main() {
	cfg := config.LoadConfig()
	store := storage.NewStorage()
	server := api.NewServer(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", server.CreateEventHandler)
	mux.HandleFunc("/update_event", server.UpdateEventHandler)
	mux.HandleFunc("/delete_event", server.DeleteEventHandler)
	mux.HandleFunc("/events_for_day", server.EventsForDayHandler)
	mux.HandleFunc("/events_for_week", server.EventsForWeekHandler)
	mux.HandleFunc("/events_for_month", server.EventsForMonthHandler)

	loggedMux := api.LoggingMiddleware(mux)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), loggedMux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
