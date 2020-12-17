package main

import (
	"context"
	"flag"
	oslog "log"
	"os"
	"os/signal"
	"time"

	internalserver "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/server"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	config "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	logger "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/logger"
	store "github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf, err := config.NewConfig(configFile)
	if err != nil {
		oslog.Fatal(err)
		return
	}
	_, err = logger.Init(conf)
	if err != nil {
		oslog.Fatal(err)
		return
	}

	storage := store.NewStore(conf)
	var calendar = app.New(storage)

	server := internalserver.NewServer(conf, calendar)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Log.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Log.Error("failed to start http server: " + err.Error())
		os.Exit(1) //nolint
	}
}
