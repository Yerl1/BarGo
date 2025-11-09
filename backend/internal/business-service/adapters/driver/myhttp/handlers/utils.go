package handlers

import (
	"encoding/json"
	"net/http"
)

// JsonResponse writes the given data as a JSON-encoded HTTP response with status code 200 OK.
func JsonResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}
