package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

type Client struct {
	address     string
	timeout     time.Duration
	conn        net.Conn
	in          io.ReadCloser
	out         io.Writer
	connScanner *bufio.Scanner
	inScanner   *bufio.Scanner
	stop        bool
}

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

func NewTelnetClient(
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *Client) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)

	if err != nil {
		return err //nolint:wrapcheck
	}

	t.connScanner = bufio.NewScanner(t.conn)
	t.inScanner = bufio.NewScanner(t.in)

	return nil
}

func (t *Client) Receive() (err error) {
	if t.conn == nil {
		return
	}
	for t.connScanner.Scan() {
		text := t.connScanner.Text()
		_, err = t.out.Write([]byte(text + "\n"))
		if err != nil {
			return err //nolint:wrapcheck
		}
	}

	if !t.stop {
		t.stop = true
		log.Println("...Remote Server stopped")
	}

	return nil
}

func (t *Client) Send() (err error) {
	if t.conn == nil {
		return
	}
	for t.inScanner.Scan() {
		text := t.inScanner.Text()

		_, err = t.conn.Write([]byte(text + "\n"))
		if err != nil {
			return err //nolint:wrapcheck
		}
	}

	if !t.stop {
		t.stop = true
		log.Println("...EOF")
	}

	return nil
}
func (t *Client) Close() (err error) {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
