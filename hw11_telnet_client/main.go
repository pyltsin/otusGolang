package main

import (
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
	timeout time.Duration
)

const (
	minLenArgs = 3
	maxLenArgs = 4
)

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout=time")
}

func main() {
	flag.Parse()
	args := os.Args

	if (len(args) < minLenArgs) || (len(args) > maxLenArgs) {
		log.Fatal(errors.New("not enough arguments, should be 3 at least"))
	}
	host := os.Args[len(os.Args)-2]
	port := os.Args[len(os.Args)-1]

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

	errs := make(chan error, 1)
	go func() { errs <- c.Send() }()
	go func() { errs <- c.Receive() }()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case <-signals:
		signal.Stop(signals)
		return

	case err = <-errs:
		if err != nil {
			log.Panic(err)
		}
		return
	}
}
