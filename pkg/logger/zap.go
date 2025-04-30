package logger

import (
	"fmt"
	"ml-prediction/config"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(cfg config.Configuration) (*zap.Logger, error) {
	if cfg.App.Environment == "production" {
		logger, err := zap.NewProduction(
			zap.Hooks(func(e zapcore.Entry) error {
				if e.Level == zapcore.ErrorLevel {
					defer sentry.Flush(2 * time.Second)
					sentry.CaptureMessage(fmt.Sprintf("%s, Line No: %d :: %s", e.Caller.File, e.Caller.Line, e.Message))
				}
				return nil
			}),
		)
		if err != nil {
			return nil, err
		}

		return logger, nil
	} else {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		return logger, nil
	}
}
