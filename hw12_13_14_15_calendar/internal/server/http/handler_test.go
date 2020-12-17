package internalhttp

import (
	"bytes"
	"encoding/json"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EventHandlerSuite struct {
	suite.Suite
	handler *GeneralHandler
	storage *memory.Storage
}

func (s *EventHandlerSuite) SetupTest() {
	storage := memory.New()
	s.handler = NewEventHandler(app.New(storage))
	s.storage = storage
}

func TestEventHandlerSuite(t *testing.T) {
	suite.Run(t, new(EventHandlerSuite))
}

func (s *EventHandlerSuite) TestGetEvent() {
	title := "test"
	id, _ := s.storage.Create(app.Event{Title: title})

	req := httptest.NewRequest(http.MethodGet, "/event/"+string(id), bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code)

	trim := strings.Trim(rec.Body.String(), "\n")
	s.Require().True(strings.Contains(trim, title))
}

func (s *EventHandlerSuite) TestGetEventList() {
	title1 := "test1"
	title2 := "test2"
	_, _ = s.storage.Create(app.Event{Title: title1})
	_, _ = s.storage.Create(app.Event{Title: title2})

	req := httptest.NewRequest(http.MethodGet, "/event", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code)

	trim := strings.Trim(rec.Body.String(), "\n")
	s.Require().True(strings.Contains(trim, title1))
	s.Require().True(strings.Contains(trim, title2))
}

func (s *EventHandlerSuite) TestUpdateEvent() {

	id, _ := s.storage.Create(app.Event{Title: "test"})

	newTitle := "new test"
	reqEvent := &app.Event{
		Title: newTitle,
	}

	body, err := json.Marshal(reqEvent)
	if err != nil {
		s.T().Fail()
	}
	req := httptest.NewRequest(http.MethodPatch, "/event/"+string(id), bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)

	event, _ := s.storage.GetByID(id)
	s.Require().Equal(newTitle, event.Title)
}

func (s *EventHandlerSuite) TestCreateEvent() {

	reqEvent := &app.Event{
		Title: "new event",
	}

	body, err := json.Marshal(reqEvent)
	if err != nil {
		s.T().Fail()
	}

	req := httptest.NewRequest(http.MethodPatch, "/event/12", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)

	list, err := s.storage.List()
	s.Require().True(len(list) == 1)
}

func (s *EventHandlerSuite) TestDeleteProject() {

	id, _ := s.storage.Create(app.Event{Title: "test"})

	newTitle := "new test"
	reqEvent := &app.Event{
		Title: newTitle,
	}

	body, err := json.Marshal(reqEvent)
	if err != nil {
		s.T().Fail()
	}
	req := httptest.NewRequest(http.MethodDelete, "/event/"+string(id), bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)

	_, found := s.storage.GetByID(id)
	s.Require().False(found)
}
