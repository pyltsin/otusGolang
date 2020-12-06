package app

import "time"

type EventID string

type Event struct {
	ID      EventID
	Title   string
	Date    time.Time
	Latency time.Duration
	Note    string
	UserID  int64
	Notify  time.Duration
}
