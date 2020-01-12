package ctx

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type key string

func (k key) String() string {
	return string(k)
}

// Context keys.
const (
	Logger key = "log"
)

// WithLog ...
func WithLog(c context.Context, l *zap.Logger) context.Context {
	return context.WithValue(c, Logger, l)
}

// WithLogReq ...
func WithLogReq(req *http.Request, l *zap.Logger) context.Context {
	return WithLog(req.Context(), l)
}

// GetLog ...
func GetLog(c context.Context) *zap.Logger {
	v := c.Value(Logger)
	if l, ok := v.(*zap.Logger); ok {
		return l
	}
	return nil
}

// GetReqLog ...
func GetReqLog(req *http.Request) *zap.Logger {
	return GetLog(req.Context())
}
