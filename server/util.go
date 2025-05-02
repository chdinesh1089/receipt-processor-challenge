package server

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const RequestIDKey ContextKey = "request_id"

func getLogger(ctx context.Context) *logrus.Entry {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	
	if request_id, ok := ctx.Value(RequestIDKey).(string); ok {
        return logger.WithField("request_id", request_id)
    }
	return logrus.NewEntry(logger)
}
