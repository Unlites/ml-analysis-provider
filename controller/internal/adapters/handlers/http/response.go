package handler

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPHandler) encodeResponse(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
