package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
