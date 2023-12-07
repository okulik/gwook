package web

import (
	"encoding/json"
	"net/http"
)

type ErrorWithCode struct {
	Err  error
	Code int
}

func NewErrorWithCode(err error, code int) *ErrorWithCode {
	return &ErrorWithCode{
		Err:  err,
		Code: code,
	}
}

func WriteErrorResponse(w http.ResponseWriter, err error, code int) *ErrorWithCode {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		return NewErrorWithCode(err, http.StatusInternalServerError)
	}
	return nil
}
