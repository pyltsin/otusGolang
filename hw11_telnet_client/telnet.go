package main

import (
	"bufio"
	"errors"
	"io"
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
		return err
	}

	t.connScanner = bufio.NewScanner(t.conn)
	t.inScanner = bufio.NewScanner(t.in)

	return nil
}

func (t *Client) Receive() (err error) {
	if t.conn == nil {
		return
	}
	if !t.connScanner.Scan() {
		return errors.New("...Connection was closed by peer")
	}
	_, err = t.out.Write([]byte(t.connScanner.Text() + "\n"))
	return err
}

func (t *Client) Send() (err error) {
	if t.conn == nil {
		return
	}
	if !t.inScanner.Scan() {
		return errors.New("...EOF")
	}
	_, err = t.conn.Write([]byte(t.inScanner.Text() + "\n"))
	return err
}

func (t *Client) Close() (err error) {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
