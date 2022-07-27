package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger = baseLogger.Sugar()
}

func Log() *zap.SugaredLogger {
	return logger
}
