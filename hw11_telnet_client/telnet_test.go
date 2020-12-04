package main

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("connect to wrong host", func(t *testing.T) {
		d, err := time.ParseDuration("2s")
		require.NoError(t, err)
		client := NewTelnetClient("test.tu.ru", d, os.Stdin, os.Stdout)
		require.Error(t, client.Connect())
	})

	t.Run("client shutdown", func(t *testing.T) {
		l, err := net.Listen("tcp", "localhost:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			_, err = conn.Read(request)
			require.NoError(t, err)

			_, err = conn.Read(request)
			require.EqualError(t, err, io.EOF.Error())
		}()

		wg.Wait()
	})

	t.Run("server down", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			closer := ioutil.NopCloser(in)
			client := NewTelnetClient(l.Addr().String(), timeout, closer, out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			_ = closer.Close()
			err = client.Send()
			require.NoError(t, err)

		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			_, err = conn.Read(request)
			assert.True(t, errors.Is(err, io.EOF))
		}()

		wg.Wait()
	})
}
