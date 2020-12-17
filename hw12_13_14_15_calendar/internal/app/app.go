package app

import (
	"context"
	"errors"
	"time"
)

type Application interface {
	CreateEvent(ctx context.Context, title string, date time.Time, latency time.Duration, note string, userID int64, notify time.Duration) (EventID, error)
	GetEvent(ctx context.Context, id EventID) (Event, error)
	GetEventList(ctx context.Context) ([]Event, error)
	UpdateEvent(ctx context.Context, id EventID, title string, date time.Time, latency time.Duration, note string, userID int64, notify time.Duration) (Event, error)
	DeleteEvent(ctx context.Context, id EventID) error
}

type App struct {
	store Store
}

var (
	errContextCanceled = errors.New("context canceled")
	ErrNotFound        = errors.New("not found")
)

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

func (a *App) UpdateEvent(ctx context.Context, id EventID, title string, date time.Time, latency time.Duration, note string, userID int64, notify time.Duration) (Event, error) {
	select {
	case <-ctx.Done():
		return Event{}, errContextCanceled
	default:
		event := Event{
			ID:      id,
			Title:   title,
			Date:    date,
			Latency: latency,
			Note:    note,
			UserID:  userID,
			Notify:  notify,
		}
		err := a.store.Update(
			event,
		)
		return event, err //nolint:wrapcheck
	}
}
func (a *App) CreateEvent(ctx context.Context, title string, date time.Time, latency time.Duration, note string, userID int64, notify time.Duration) (EventID, error) {
	select {
	case <-ctx.Done():
		return "", errContextCanceled
	default:
		return a.store.Create(
			Event{ //nolint:exhaustivestruct
				Title:   title,
				Date:    date,
				Latency: latency,
				Note:    note,
				UserID:  userID,
				Notify:  notify,
			},
		)
	}
}

func (a *App) GetEvent(ctx context.Context, id EventID) (Event, error) {
	select {
	case <-ctx.Done():
		return Event{}, errContextCanceled
	default:
		byID, find := a.store.GetByID(id)
		if find {
			return byID, nil
		}
		return Event{}, ErrNotFound
	}
}

func (a *App) GetEventList(ctx context.Context) ([]Event, error) {
	select {
	case <-ctx.Done():
		return nil, errContextCanceled
	default:
		return a.store.List()
	}
}

func (a *App) DeleteEvent(ctx context.Context, id EventID) error {
	select {
	case <-ctx.Done():
		return errContextCanceled
	default:
		return a.store.Delete(id)
	}
}
