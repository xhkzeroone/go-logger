package logger

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

func GetRequestID(ctx context.Context) string {
	v := ctx.Value("requestId")
	if v == nil {
		return "null"
	}
	return fmt.Sprint(v)
}

func WithContext(ctx context.Context) *logrus.Entry {
	requestID := GetRequestID(ctx)
	return Log.WithFields(logrus.Fields{
		"requestId": requestID,
	})
}
