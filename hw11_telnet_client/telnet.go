package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

type TelnetClient interface {
	Connect() error
	Receive() error
	Send() error
	Close() error
}

func (c *Client) Connect() (err error) {
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	return fmt.Errorf("not connect: %w", err)
}

func (c *Client) Send() error {
	_, err := io.Copy(c.conn, c.in)
	return fmt.Errorf("not send: %w", err)
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	return fmt.Errorf("not receive: %w", err)
}
func (c *Client) Close() error {
	return c.conn.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
