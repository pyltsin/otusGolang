package internalhttp

import (
	"context"
	"errors"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
)

type GeneralHandler struct {
	app app.Application
}

func NewEventHandler(a app.Application) *GeneralHandler {
	return &GeneralHandler{app: a}
}

func getPathAndID(urlPath string) (string, string) {
	var path string

	lastInd := strings.LastIndex(urlPath, "/")
	if lastInd == 0 {
		return urlPath, ""
	}
	itemID := urlPath[lastInd+1:]
	if itemID != "" {
		path = urlPath[:lastInd]
	} else {
		path = urlPath
	}

	return path, itemID
}

const timeout = time.Second

func (h GeneralHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	path, id := getPathAndID(r.URL.Path)
	if path != "/event" {
		http.NotFound(w, r)

		return
	}
	if id != "" {
		switch r.Method {
		case "GET":
			h.GetEvent(ctx, w, r)
		case "PATCH":
			h.UpdateEvent(ctx, w, r)
		case "DELETE":
			h.DeleteEvent(ctx, w, r)
		default:
		}
	} else {
		switch r.Method {
		case "GET":
			h.GetEventList(ctx, w, r)
		case "POST":
			h.CreateEvent(ctx, w, r)
		default:
		}
	}
}

func (h *GeneralHandler) GetEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, id := getPathAndID(r.URL.Path)

	if checkHeader(r, w) {
		return
	}

	e, err := h.app.GetEvent(ctx, app.EventID(id))
	if handleError(err, w) {
		return
	}
	writeResponse(e, w)
}

func (h *GeneralHandler) GetEventList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if checkHeader(r, w) {
		return
	}

	e, err := h.app.GetEventList(ctx)
	if handleError(err, w) {
		return
	}
	writeResponse(e, w)
}

func (h *GeneralHandler) CreateEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if checkHeader(r, w) {
		return
	}
	e := &app.Event{}
	err := decodeBody(r, e, w)
	if err != nil {
		_ = handleError(err, w)
	}
	id, err := h.app.CreateEvent(ctx,
		e.Title,
		e.Date,
		e.Latency,
		e.Note,
		e.UserID,
		e.Notify,
	)
	if handleError(err, w) {
		return
	}
	writeResponse(id, w)
}
func (h *GeneralHandler) UpdateEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, id := getPathAndID(r.URL.Path)

	if checkHeader(r, w) {
		return
	}
	e := &app.Event{}
	err := decodeBody(r, e, w)
	if err != nil {
		_ = handleError(err, w)
	}
	event, err := h.app.UpdateEvent(ctx,
		app.EventID(id),
		e.Title,
		e.Date,
		e.Latency,
		e.Note,
		e.UserID,
		e.Notify,
	)
	if handleError(err, w) {
		return
	}
	writeResponse(event, w)
}

// DeleteEvent deletes event.
func (h *GeneralHandler) DeleteEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	_, id := getPathAndID(r.URL.Path)

	err := h.app.DeleteEvent(ctx, app.EventID(id))
	if handleError(err, w) {
		return
	}
	writeResponse(nil, w)
}

func decodeBody(r *http.Request, e interface{}, w http.ResponseWriter) error {
	if e == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return nil
	}
	return err
}

func writeResponse(e interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func handleError(err error, w http.ResponseWriter) bool {
	switch {
	case err == nil:
	case errors.Is(err, app.ErrNotFound):
		http.Error(w, "event not found", http.StatusNotFound)
		return true
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func checkHeader(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)

		return true
	}
	return false
}
