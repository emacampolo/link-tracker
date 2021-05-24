package handler

import (
	"errors"
	"net/http"
	"strconv"

	"link-tracker/internal/link"
	"link-tracker/internal/platform/web"
)

type Link struct {
	linkService link.Service
}

func NewLink(l link.Service) *Link {
	return &Link{
		linkService: l,
	}
}

func (l *Link) Create() web.Handler {
	type request struct {
		Link     string `json:"link"`
		Password string `json:"password"`
	}

	type response struct {
		ID int `json:"id"`
	}

	return func(w http.ResponseWriter, req *http.Request) error {
		var r request
		if err := web.Decode(req, &r); err != nil {
			return web.NewError(400, err.Error())
		}

		l, err := l.linkService.Create(req.Context(), r.Link, r.Password)
		if err != nil {
			return err
		}

		resp := response{
			ID: l.ID,
		}

		return web.Respond(req.Context(), w, resp, http.StatusCreated)
	}
}

func (l *Link) Redirect() web.Handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		idParam := web.Param(req, "id")
		if idParam == "" {
			return web.NewError(400, "id param is missing")
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			return web.NewError(400, err.Error())
		}

		password := req.URL.Query().Get("password")
		if password == "" {
			return web.NewError(400, "password is missing")
		}

		ll, err := l.linkService.Redirect(req.Context(), id, password)
		if err != nil {
			if errors.Is(err, link.ErrNotFound) {
				return web.NewError(404, err.Error())
			}

			return err
		}

		http.Redirect(w, req, ll.URL, http.StatusMovedPermanently)
		return nil
	}
}
