package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteError(w http.ResponseWriter, code int, err error, msg string) {
	setHeader(w, code)
	json.NewEncoder(w).Encode(map[string]string{msg: err.Error()})
}

func WriteResponse(w http.ResponseWriter, code int, val any) {
	setHeader(w, code)
	json.NewEncoder(w).Encode(val)
}

func setHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
