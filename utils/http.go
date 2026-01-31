package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ReturnJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetIdFromRequest(w http.ResponseWriter, r *http.Request, path string) string {
	id := strings.TrimPrefix(r.URL.Path, path)
	return id
}
