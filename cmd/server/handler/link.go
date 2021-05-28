package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/emacampolo/link-tracker/internal/platform/web"
)

type Link struct {
	linkService link.Service
}

func NewLink(l link.Service) *Link {
	return &Link{
		linkService: l,
	}
}

func (lnk *Link) Create() web.Handler {
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
			return web.NewError(http.StatusBadRequest, err.Error())
		}

		if r.Link == "" {
			return web.NewError(http.StatusBadRequest, "link is missing")
		}

		if r.Password == "" {
			return web.NewError(http.StatusBadRequest, "password is missing")
		}

		l, err := lnk.linkService.Create(req.Context(), r.Link, r.Password)
		if err != nil {
			return err
		}

		resp := response{
			ID: l.ID,
		}

		return web.Respond(req.Context(), w, resp, http.StatusCreated)
	}
}

func (lnk *Link) Redirect() web.Handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		id, err := lnk.extractID(req)
		if err != nil {
			return web.NewError(http.StatusBadRequest, err.Error())
		}

		password := req.URL.Query().Get("password")
		if password == "" {
			return web.NewError(http.StatusBadRequest, "password is missing")
		}

		ll, err := lnk.linkService.Redirect(req.Context(), id, password)
		if err != nil {
			if errors.Is(err, link.ErrNotFound) {
				return web.NewError(http.StatusNotFound, err.Error())
			}

			if errors.Is(err, link.ErrInactive) {
				return web.NewError(http.StatusUnprocessableEntity, err.Error())
			}

			return err
		}

		http.Redirect(w, req, ll.URL, http.StatusMovedPermanently)
		return nil
	}
}

func (lnk *Link) Metrics() web.Handler {
	type response struct {
		ID       int    `json:"id"`
		URL      string `json:"url"`
		Count    int    `json:"count"`
		Inactive bool   `json:"inactive"`
	}

	return func(w http.ResponseWriter, req *http.Request) error {
		id, err := lnk.extractID(req)
		if err != nil {
			return web.NewError(http.StatusBadRequest, err.Error())
		}

		l, err := lnk.linkService.FindByID(req.Context(), id)
		if err != nil {
			if errors.Is(err, link.ErrNotFound) {
				return web.NewError(http.StatusNotFound, err.Error())
			}

			return err
		}

		resp := response{
			ID:       l.ID,
			URL:      l.URL,
			Count:    l.Count,
			Inactive: l.Inactive,
		}

		return web.Respond(req.Context(), w, resp, http.StatusOK)
	}
}

func (lnk *Link) Inactivate() web.Handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		id, err := lnk.extractID(req)
		if err != nil {
			return web.NewError(http.StatusBadRequest, err.Error())
		}

		if err := lnk.linkService.Inactivate(req.Context(), id); err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func (lnk *Link) extractID(req *http.Request) (int, error) {
	idParam := web.Param(req, "id")
	if idParam == "" {
		return 0, web.NewError(http.StatusBadRequest, "id param is missing")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, web.NewError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}
