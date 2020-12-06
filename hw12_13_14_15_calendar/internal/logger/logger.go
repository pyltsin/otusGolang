package logger

import (
	"errors"
	oslog "log"
	"strings"
	"time"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLogPath = "log/log.out"
)

type Logger struct {
	logger *zap.Logger
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
	_ = l.logger.Sync()
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
	_ = l.logger.Sync()
}

var (
	Log *Logger
)

func Init(config config.Config) (*Logger, error) {
	var lvl zap.AtomicLevel
	switch strings.ToUpper(config.Logger.Level) {
	case "DEBUG":
		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "ERROR":
		lvl = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "INFO":
		lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "WARN":
		lvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		return nil, errors.New("not found log level")
	}

	var file string
	if config.Logger.File == "" {
		file = DefaultLogPath
	} else {
		file = config.Logger.File
	}

	log, err := zap.Config{ //nolint:exhaustivestruct
		Level:       lvl,
		Encoding:    "console",
		OutputPaths: []string{file},
		EncoderConfig: zapcore.EncoderConfig{ //nolint:exhaustivestruct
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.Stamp))
			},
		},
	}.Build()

	if err != nil {
		oslog.Fatal(err)
	}

	Log = &Logger{
		logger: log,
	}
	return Log, nil
}
