package storage

import (
	"errors"
	"sync"
	"time"

	"calendar/internal/app/models"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("date is busy")
)

type Storage struct {
	mu     sync.RWMutex
	events map[int]models.Event
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		events: make(map[int]models.Event),
		nextID: 1,
	}
}


func (s *Storage) CreateEvent(event models.Event) (models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем занятость даты
	for _, e := range s.events {
		if e.UserID == event.UserID &&
			e.Date.Year() == event.Date.Year() &&
			e.Date.Month() == event.Date.Month() &&
			e.Date.Day() == event.Date.Day() {
			return models.Event{}, ErrDateBusy
		}
	}

	event.ID = s.nextID
	s.events[event.ID] = event
	s.nextID++
	return event, nil
}


func (s *Storage) UpdateEvent(event models.Event) (models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return models.Event{}, ErrEventNotFound
	}

	s.events[event.ID] = event
	return event, nil
}

func (s *Storage) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) EventsForDay(userID int, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Event
	for _, event := range s.events {
		if event.UserID == userID && event.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) EventsForWeek(userID int, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Event
	year, week := date.ISOWeek()
	for _, event := range s.events {
		eYear, eWeek := event.Date.ISOWeek()
		if event.UserID == userID && eYear == year && eWeek == week {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) EventsForMonth(userID int, date time.Time) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Event
	for _, event := range s.events {
		if event.UserID == userID && event.Date.Month() == date.Month() && event.Date.Year() == date.Year() {
			result = append(result, event)
		}
	}
	return result, nil
}
