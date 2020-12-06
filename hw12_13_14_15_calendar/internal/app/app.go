package app

import (
	"context"
)

type App struct {
	store Store
}

type Store interface {
	Create(event Event) (EventID, error)
	Update(event Event) error
	Delete(id EventID) error
	List() ([]Event, error)
	GetByID(id EventID) (Event, bool)
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func New(store Store) *App {
	return &App{store: store}
}

func (a *App) CreateEvent(ctx context.Context, title string) (EventID, error) {
	return a.store.Create(
		Event{ //nolint:exhaustivestruct
			Title: title,
		},
	)
}
