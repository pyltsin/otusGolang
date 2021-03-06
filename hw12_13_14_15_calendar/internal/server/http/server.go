package internalhttp

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	conf   config.ServerConf
	app    *Application
	server *http.Server
}

type Application interface {
	CreateEvent(ctx context.Context, title string) (app.EventID, error)
	// TODO
}

func NewServer(conf config.Config, app Application) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(HelloWorld)))

	server := &http.Server{ //nolint:exhaustivestruct
		Addr:         net.JoinHostPort(conf.Server.Address, conf.Server.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		conf:   conf.Server,
		app:    &app,
		server: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error("Could not listen")
		return err //nolint:wrapcheck
	}

	if <-ctx.Done(); true {
		return s.Stop(ctx)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Log.Error("Could not gracefully shutdown the server")
		return err //nolint:wrapcheck
	}
	return nil
}

// HelloWorld simple route.
func HelloWorld(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "hello-world")
}
