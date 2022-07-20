package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, code int, data map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(data)
	w.WriteHeader(code)
	w.Write(resp)
}
