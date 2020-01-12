package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Write writes data to response writer.
func Write(c context.Context, w http.ResponseWriter, data interface{}) *AppError {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return &AppError{
			Err:     err,
			Message: fmt.Sprintf("error in response write: %v", err),
			Code:    http.StatusInternalServerError,
			Context: c,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	return nil
}
