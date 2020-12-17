package internalserver

import (
	"context"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	internalgrpc "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/server/http"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func NewServer(conf config.Config, app app.Application) Server {
	if conf.Grpc.Enable {
		return internalgrpc.NewServer(conf, app)
	}
	return internalhttp.NewServer(conf, app)
}
