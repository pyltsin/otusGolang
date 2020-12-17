package logger

import (
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "log.")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	conf := config.Config{Logger: config.LoggerConf{
		File: tmpfile.Name(), Level: "warn"}}

	log, err := Init(conf)
	require.NoError(t, err)

	t.Run("Test message", func(t *testing.T) {
		log.Info("debug message")
		log.Error("error message")

		res, err := ioutil.ReadAll(tmpfile)
		require.NoError(t, err)

		require.Less(t, strings.Index(string(res), "debug message"), 0)
		require.Greater(t, strings.Index(string(res), "error message"), 0)
	})

}
