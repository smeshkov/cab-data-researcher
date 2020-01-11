package handlers

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// AppHandler ...
// http://blog.golang.org/error-handling-and-go
type AppHandler func(http.ResponseWriter, *http.Request) *AppError

// AppError ...
type AppError struct {
	Err     error
	Message string
	Code    int
	Log     *zap.Logger
}

func (e *AppError) Error() string {
	return fmt.Sprintf("application error code %d: %s - %v", e.Code, e.Message, e.Err)
}

// ServeHTTP ...
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		l := e.Log
		if l == nil {
			l = zap.L()
		}

		l.Error("handler error",
			zap.Int("status_code", e.Code),
			zap.String("err_message", e.Message),
			zap.Error(e.Err))

		http.Error(w, e.Message, e.Code)
	}
}

// AppErrorf ...
func AppErrorf(err error, format string, v ...interface{}) *AppError {
	return &AppError{
		Err:     err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
