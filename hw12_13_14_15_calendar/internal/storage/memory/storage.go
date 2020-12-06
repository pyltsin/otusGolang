package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	Events map[app.EventID]app.Event
	Mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{Events: make(map[app.EventID]app.Event)} //nolint:exhaustivestruct
}

func (s *Storage) Create(event app.Event) (app.EventID, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var id = uuid.New()
	event.ID = app.EventID(id.String())

	s.Events[event.ID] = event

	return event.ID, nil
}

func (s *Storage) Update(event app.Event) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id app.EventID) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Events, id)
	return nil
}

func (s *Storage) List() ([]app.Event, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	result := make([]app.Event, 0, len(s.Events))

	for _, value := range s.Events {
		result = append(result, value)
	}

	return result, nil
}

func (s *Storage) GetByID(id app.EventID) (app.Event, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	event, ok := s.Events[id]
	if !ok {
		return app.Event{}, false
	}
	return event, true
}
