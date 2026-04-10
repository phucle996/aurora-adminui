package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// JSON writes the shared API response payload through a standard net/http writer.
func JSON(w http.ResponseWriter, status int, message string, data any) {
	if w == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Message: message,
		Data:    data,
	})
}
