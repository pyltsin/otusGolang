package sqlstorage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
)

const dateTimeLayout = "2006-01-02 15:04:00 -0700"

type Storage struct {
	db     *sql.DB
	config config.StorageConf
}

func New(conf config.StorageConf) *Storage {
	return &Storage{config: conf}
}

func (s *Storage) Connect() error {
	var err error
	s.db, err = sql.Open("mysql", s.config.SQLUser+":"+s.config.SQLPass+"@tcp("+s.config.SQLHost+":"+s.config.SQLPort+")/"+s.config.SQLDbase)
	if err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Create(ev app.Event) (app.EventID, error) {
	ev.ID = app.EventID(uuid.New().String())
	_, err := s.db.Exec(
		`INSERT INTO events 
		(id, title, date_event, latency, note, userID, notifyTime) VALUES 
		($1, $2, $3, $4, $5, $6)`,
		ev.ID,
		ev.Title,
		ev.Date.Format(dateTimeLayout),
		ev.Latency,
		ev.Note,
		ev.UserID,
		ev.Notify,
	)
	if err != nil {
		return "", err //nolint:wrapcheck
	}
	return ev.ID, err //nolint:wrapcheck
}

func (s *Storage) Update(event app.Event) error {
	_, err := s.db.Exec(
		`UPDATE events set 
		title=$1, date_event=$2, latency=$3, note=$4, userID=$5, notifyTime=$6
		where id=$7`,
		event.Title,
		event.Date.Format(dateTimeLayout),
		event.Latency,
		event.Note,
		event.UserID,
		event.Notify,
		event.ID,
	)
	return err //nolint:wrapcheck
}

func (s *Storage) Delete(id app.EventID) error {
	_, err := s.db.Exec(
		`DELETE from events where id=$1`,
		id,
	)
	return err //nolint:wrapcheck
}

func (s *Storage) List() ([]app.Event, error) {
	result := make([]app.Event, 0)
	results, err := s.db.Query(`SELECT (id,title,date,latency,note,userID,notifyTime) from events ORDER BY id`) //nolint:rowserrcheck
	if err != nil {
		return result, err //nolint:wrapcheck
	}
	defer results.Close()
	for results.Next() {
		var evt app.Event
		var dateRaw string
		err = results.Scan(&evt.ID, &evt.Title, &dateRaw, &evt.Latency, &evt.Note, &evt.UserID, &evt.Notify)
		if err != nil {
			return result, err //nolint:wrapcheck
		}
		evt.Date, err = time.Parse(dateTimeLayout, dateRaw)
		if err != nil {
			return result, err //nolint:wrapcheck
		}
		result = append(result, evt)
	}
	return result, nil
}

func (s *Storage) GetByID(id app.EventID) (app.Event, bool) {
	var res app.Event
	var dateRaw string
	err := s.db.QueryRow(
		`SELECT (id,title,date_event,latency,note,userID,notifyTime) from events where id=$1`, id).Scan(res.ID, res.Title, dateRaw, res.Latency, res.Note, res.UserID, res.Notify)
	if err != nil {
		return res, false
	}
	parsedDate, err := time.Parse(dateTimeLayout, dateRaw)
	if err != nil {
		return res, false
	}
	res.Date = parsedDate
	return res, true
}
