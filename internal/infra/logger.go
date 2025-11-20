package infra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Zap *zap.Logger
	*zap.SugaredLogger
	LogsPath   string
	Name       string
	FileWriter *os.File
}

func NewLogger(lc fx.Lifecycle) (*Logger, error) {
	var l Logger
	l.Name = "main"
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	l.LogsPath = filepath.Join(wd, "./logs")
	err = os.MkdirAll(l.LogsPath, 0o750)
	if err != nil {
		return nil, err
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		NameKey:     "logger",
		CallerKey:   "caller",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(time.FixedZone("GMT+0", 3*60*60)).Format("2006-01-02:15"))
		},
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	level := zap.InfoLevel

	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	fileEncoderConfig := encoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	var cores []zapcore.Core

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	cores = append(cores, consoleCore)

	logPath := filepath.Join(l.LogsPath, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02_15")))

	// gosec: G304 is acceptable here as logPath is constructed from safe components
	// nolint:gosec
	fileWriter, errOpen := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if errOpen != nil {
		return nil, errOpen
	}

	l.FileWriter = fileWriter

	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), level)
	cores = append(cores, fileCore)

	combinedCore := zapcore.NewTee(cores...)

	log := zap.New(combinedCore, zap.AddCaller())
	l.SugaredLogger = log.Sugar()
	l.Zap = log
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			l.Info("logger initialized")

			return nil
		},
		OnStop: func(_ context.Context) error {
			if l.FileWriter != nil {
				return l.FileWriter.Close()
			}
			l.Info("logger stopped")
			return l.Sync()
		},
	})

	return &l, nil
}

type ZapFxLogger struct{ *zap.Logger }

func (z *ZapFxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		z.Info("fx on start executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName))
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			z.Error("fx on start failed",
				zap.String("callee", e.FunctionName),
				zap.Error(e.Err))
		} else {
			z.Info("fx on start succeeded",
				zap.String("callee", e.FunctionName),
				zap.Duration("runtime", e.Runtime))
		}
	default:
		z.Debug("fx event", zap.String("type", fmt.Sprintf("%T", e)))
	}
}

// type ZapGooseAdapter struct{ zap *zap.Logger }

// func (a *ZapGooseAdapter) Printf(format string, args ...any) {
// 	a.zap.Sugar().Infof(format, args...)
// }

// func (a *ZapGooseAdapter) Fatalf(format string, args ...interface{}) {
// 	a.zap.Sugar().Fatalf(format, args...)
// }
