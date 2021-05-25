package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"link-tracker/cmd/server/handler"
	"link-tracker/internal/link"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type linkServiceMock struct {
	mock.Mock
}

func (l *linkServiceMock) Create(ctx context.Context, url, password string) (link.Link, error) {
	args := l.Called(ctx, url, password)
	return args.Get(0).(link.Link), args.Error(1)
}

func (l *linkServiceMock) Redirect(ctx context.Context, ID int, password string) (link.Link, error) {
	args := l.Called(ctx, ID, password)
	return args.Get(0).(link.Link), args.Error(1)
}

func (l *linkServiceMock) FindByID(ctx context.Context, ID int) (link.Link, error) {
	args := l.Called(ctx, ID)
	return args.Get(0).(link.Link), args.Error(1)
}

func (l *linkServiceMock) Inactivate(ctx context.Context, ID int) error {
	return l.Called(ctx, ID).Error(0)
}

func TestLink_Create(t *testing.T) {
	// Given
	r := struct {
		Link     string `json:"link"`
		Password string `json:"password"`
	}{
		Link:     "https://www.google.com",
		Password: "123456",
	}

	body, _ := json.Marshal(r)
	req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	l := link.Link{ID: 1}

	svcMock := &linkServiceMock{}
	svcMock.On("Create", req.Context(), r.Link, r.Password).Return(l, nil)

	linkHandler := handler.NewLink(svcMock)

	// When
	linkHandler.Create().ServeHTTP(rr, req)

	// Then
	require.Equal(t, http.StatusCreated, rr.Code)
	require.JSONEq(t, `{"id":1}`, rr.Body.String())
}

func TestLink_Create_RequiredFields(t *testing.T) {
	type request struct {
		Link     string `json:"link"`
		Password string `json:"password"`
	}

	tt := []struct {
		name    string
		req     request
		wantErr string
	}{
		{
			name:    "link is required",
			req:     request{Password: "1234"},
			wantErr: `{"code":"bad_request","message":"link is missing"}`,
		},
		{
			name:    "password is required",
			req:     request{Link: "https://www.google.com"},
			wantErr: `{"code":"bad_request","message":"password is missing"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			l := link.Link{ID: 1}

			svcMock := &linkServiceMock{}
			svcMock.On("Create", req.Context(), tc.req.Link, tc.req.Password).Return(l, nil)

			linkHandler := handler.NewLink(svcMock)

			// When
			linkHandler.Create().ServeHTTP(rr, req)

			// Then
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Equal(t, tc.wantErr, rr.Body.String())
		})
	}

	r := struct {
		Link     string `json:"link"`
		Password string `json:"password"`
	}{
		Link:     "https://www.google.com",
		Password: "123456",
	}

	body, _ := json.Marshal(r)
	req := httptest.NewRequest(http.MethodPost, "/link", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	l := link.Link{ID: 1}

	svcMock := &linkServiceMock{}
	svcMock.On("Create", req.Context(), r.Link, r.Password).Return(l, nil)

	linkHandler := handler.NewLink(svcMock)

	// When
	linkHandler.Create().ServeHTTP(rr, req)

	// Then
	require.Equal(t, http.StatusCreated, rr.Code)
	require.JSONEq(t, `{"id":1}`, rr.Body.String())
}
