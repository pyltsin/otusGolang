package storage

import (
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStore(conf config.Config) app.Store {
	if conf.Storage.InMemory {
		st := memorystorage.New()
		return st
	}
	return sqlstorage.New(conf.Storage)
}
