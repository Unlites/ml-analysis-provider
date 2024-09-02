package httphandler

import (
	"encoding/json"
	"net/http"
)

// encodeResponse writes given http status and data to response of request
func encodeResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
