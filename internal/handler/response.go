package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code":  strconv.Itoa(code),
	})
}
