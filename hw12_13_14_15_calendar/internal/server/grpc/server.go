package internalgrpc

import (
	"context"
	"net"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/logger"
	googrpc "google.golang.org/grpc"
)

//go:generate protoc EventService.proto --go_out=. --go-grpc_out=. -I ../../../api

type Server struct {
	conf config.GrpcConf
	app  app.Application
	s    *googrpc.Server
}

func NewServer(conf config.Config, app app.Application) *Server {
	return &Server{app: app, conf: conf.Grpc, s: googrpc.NewServer(googrpc.UnaryInterceptor(LoggingInterceptor))}
}

func (s *Server) Start(ctx context.Context) error {
	logger.Log.Info("GRPC server starting")
	listnGrpc, err := net.Listen("tcp", net.JoinHostPort(s.conf.Address, s.conf.Port))
	RegisterEventsServer(s.s, NewInternalAdapter(s.app))
	if err != nil {
		return err
	}
	return s.s.Serve(listnGrpc)
}

func (s *Server) Stop(ctx context.Context) error {
	s.s.GracefulStop()
	return nil
}
