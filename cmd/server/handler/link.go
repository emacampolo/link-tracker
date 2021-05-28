package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/gin-gonic/gin"
)

type Link struct {
	linkService link.Service
}

func NewLink(l link.Service) *Link {
	return &Link{
		linkService: l,
	}
}

func (lnk *Link) Create(c *gin.Context) {
	type request struct {
		Link     string `json:"link"`
		Password string `json:"password"`
	}

	type response struct {
		ID int `json:"id"`
	}

	var r request

	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if r.Link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link is missing"})
		return
	}

	if r.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is missing"})
		return
	}

	l, err := lnk.linkService.Create(c.Request.Context(), r.Link, r.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response{
		ID: l.ID,
	}

	c.JSON(http.StatusCreated, resp)
}

func (lnk *Link) Redirect(c *gin.Context) {
	id, err := lnk.extractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password := c.Query("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is missing"})
		return
	}

	ll, err := lnk.linkService.Redirect(c.Request.Context(), id, password)
	if err != nil {
		if errors.Is(err, link.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, link.ErrInactive) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusMovedPermanently, ll.URL)
}

func (lnk *Link) Metrics(c *gin.Context) {
	type response struct {
		ID       int    `json:"id"`
		URL      string `json:"url"`
		Count    int    `json:"count"`
		Inactive bool   `json:"inactive"`
	}

	id, err := lnk.extractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l, err := lnk.linkService.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, link.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response{
		ID:       l.ID,
		URL:      l.URL,
		Count:    l.Count,
		Inactive: l.Inactive,
	}

	c.JSON(http.StatusOK, resp)
}

func (lnk *Link) Inactivate(c *gin.Context) {
	id, err := lnk.extractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := lnk.linkService.Inactivate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (lnk *Link) extractID(c *gin.Context) (int, error) {
	idParam := c.Param("id")
	if idParam == "" {
		return 0, errors.New("id param is missing")
	}

	return strconv.Atoi(idParam)
}
