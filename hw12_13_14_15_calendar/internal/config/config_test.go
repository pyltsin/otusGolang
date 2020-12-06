package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestNewConfig(t *testing.T) {

	t.Run("Bad file", func(t *testing.T) {
		c, e := NewConfig("./../../testfiles/badfile_config.toml")
		require.Equal(t, Config{}, c)
		require.Error(t, e)
	})

	t.Run("TOML reading", func(t *testing.T) {
		config, e := NewConfig("./../../testfiles/goodfile_config.toml")
		require.Equal(t, false, config.Storage.InMemory)
		require.Equal(t, "INFO", config.Logger.Level)
		require.Equal(t, "8080", config.Server.Port)
		require.NoError(t, e)
	})

}
