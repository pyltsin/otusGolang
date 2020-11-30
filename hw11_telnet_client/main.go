package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout          time.Duration
	ErrNotEnoughArgs = errors.New("should be 3 at least")
	signals          = make(chan os.Signal, 1)
)

const (
	minLenArgs = 3
)

func init() {
	flag.DurationVar(&timeout, "timeout", 0, "connection timeout")
}

func main() {
	flag.Parse()
	if len(os.Args) < minLenArgs {
		log.Fatal(ErrNotEnoughArgs)
	}
	host := os.Args[2]
	port := os.Args[3]

	c := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go work(c.Receive, cancelFunc)
	go work(c.Send, cancelFunc)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signals:
		cancelFunc()
		signal.Stop(signals)
		return

	case <-ctx.Done():
		close(signals)
		return
	}
}

func work(handler func() error, cancel context.CancelFunc) {
	if err := handler(); err != nil {
		cancel()
	}
}
