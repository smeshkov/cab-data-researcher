package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Write writes data to response writer.
func Write(rw http.ResponseWriter, data interface{}) *AppError {
	return WriteWithLog(rw, data, zap.L())
}

// WriteWithLog writes data to response writer.
func WriteWithLog(rw http.ResponseWriter, data interface{}, l *zap.Logger) *AppError {
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		return &AppError{
			Err:     err,
			Message: fmt.Sprintf("error in response write: %v", err),
			Code:    http.StatusInternalServerError,
			Log:     l,
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	return nil
}
