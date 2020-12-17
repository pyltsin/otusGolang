package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

type LoggerConf struct {
	File  string
	Level string
}

type StorageConf struct {
	InMemory bool
	SQLHost  string
	SQLPort  string
	SQLDbase string
	SQLUser  string
	SQLPass  string
}

type ServerConf struct {
	Address string
	Port    string
}

func NewConfig(configFile string) (Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, err //nolint:wrapcheck
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, err //nolint:wrapcheck
	}
	var config Config
	_, err = toml.Decode(string(s), &config)
	return config, err //nolint:wrapcheck
}
