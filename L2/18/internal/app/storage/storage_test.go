package storage

import (
	"testing"
	"time"

	"calendar/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage_CreateEvent(t *testing.T) {
	s := NewStorage()
	event := models.Event{
		UserID: 1,
		Date:   time.Now(),
		Event:  "Test Event",
	}
	createdEvent, err := s.CreateEvent(event)
	assert.NoError(t, err)
	assert.Equal(t, 1, createdEvent.ID)
	assert.Equal(t, event.UserID, createdEvent.UserID)
	assert.Equal(t, event.Event, createdEvent.Event)
}

func TestStorage_CreateEvent_DateBusy(t *testing.T) {
	s := NewStorage()
	date := time.Now()
	event1 := models.Event{
		UserID: 1,
		Date:   date,
		Event:  "Test Event 1",
	}
	_, err := s.CreateEvent(event1)
	assert.NoError(t, err)

	event2 := models.Event{
		UserID: 1,
		Date:   date,
		Event:  "Test Event 2",
	}
	_, err = s.CreateEvent(event2)
	assert.Error(t, err)
	assert.Equal(t, ErrDateBusy, err)
}

func TestStorage_UpdateEvent(t *testing.T) {
	s := NewStorage()
	event := models.Event{
		UserID: 1,
		Date:   time.Now(),
		Event:  "Test Event",
	}
	createdEvent, _ := s.CreateEvent(event)

	updatedEventData := models.Event{
		ID:     createdEvent.ID,
		UserID: 1,
		Date:   createdEvent.Date,
		Event:  "Updated Test Event",
	}
	updatedEvent, err := s.UpdateEvent(updatedEventData)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test Event", updatedEvent.Event)
}

func TestStorage_UpdateEvent_NotFound(t *testing.T) {
	s := NewStorage()
	event := models.Event{
		ID:     999,
		UserID: 1,
		Date:   time.Now(),
		Event:  "Test Event",
	}
	_, err := s.UpdateEvent(event)
	assert.Error(t, err)
	assert.Equal(t, ErrEventNotFound, err)
}

func TestStorage_DeleteEvent(t *testing.T) {
	s := NewStorage()
	event := models.Event{
		UserID: 1,
		Date:   time.Now(),
		Event:  "Test Event",
	}
	createdEvent, _ := s.CreateEvent(event)

	err := s.DeleteEvent(createdEvent.ID)
	assert.NoError(t, err)

	_, err = s.EventsForDay(1, createdEvent.Date)
	assert.NoError(t, err)
}

func TestStorage_DeleteEvent_NotFound(t *testing.T) {
	s := NewStorage()
	err := s.DeleteEvent(999)
	assert.Error(t, err)
	assert.Equal(t, ErrEventNotFound, err)
}

func TestStorage_EventsForDay(t *testing.T) {
	s := NewStorage()
	date := time.Now()
	event1 := models.Event{
		UserID: 1,
		Date:   date,
		Event:  "Test Event 1",
	}
	event2 := models.Event{
		UserID: 1,
		Date:   date.Add(24 * time.Hour),
		Event:  "Test Event 2",
	}
	s.CreateEvent(event1)
	s.CreateEvent(event2)

	events, err := s.EventsForDay(1, date)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, event1.Event, events[0].Event)
}

func TestStorage_EventsForWeek(t *testing.T) {
	s := NewStorage()
	date := time.Now()
	event1 := models.Event{
		UserID: 1,
		Date:   date,
		Event:  "Test Event 1",
	}
	event2 := models.Event{
		UserID: 1,
		Date:   date.Add(7 * 24 * time.Hour),
		Event:  "Test Event 2",
	}
	s.CreateEvent(event1)
	s.CreateEvent(event2)

	events, err := s.EventsForWeek(1, date)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, event1.Event, events[0].Event)
}

func TestStorage_EventsForMonth(t *testing.T) {
	s := NewStorage()
	date := time.Now()
	event1 := models.Event{
		UserID: 1,
		Date:   date,
		Event:  "Test Event 1",
	}
	event2 := models.Event{
		UserID: 1,
		Date:   date.Add(30 * 24 * time.Hour),
		Event:  "Test Event 2",
	}
	s.CreateEvent(event1)
	s.CreateEvent(event2)

	events, err := s.EventsForMonth(1, date)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, event1.Event, events[0].Event)
}
