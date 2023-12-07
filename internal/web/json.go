package web

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, v any, code int) *ErrorWithCode {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return NewErrorWithCode(err, http.StatusInternalServerError)
	}
	return nil
}
