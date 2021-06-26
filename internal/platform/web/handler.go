package web

import (
	"errors"
	"log"
	"net/http"
)

// Handler that allow default error handling.
type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// If the error was of the type *Error, the handler has a specific status code and error to return.
		var webErr *Error
		if !errors.As(err, &webErr) {
			webErr = NewErrorf(500, err.Error()).(*Error)
		}

		if err := Respond(r.Context(), w, webErr, webErr.Status); err != nil {
			log.Printf("writing http response : %v", err)
		}
	}
}
