package link_test

import (
	"context"
	"testing"

	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Save(ctx context.Context, l link.Link) int {
	return r.Mock.Called(ctx, l).Int(0)
}

func (r *repositoryMock) Update(ctx context.Context, l link.Link) error {
	return r.Mock.Called(ctx, l).Error(0)
}

func (r *repositoryMock) FindByID(ctx context.Context, ID int) (link.Link, error) {
	args := r.Mock.Called(ctx, ID)
	return args.Get(0).(link.Link), args.Error(1)
}

func TestService_Create(t *testing.T) {
	// Given
	ctx := context.Background()
	url := "https://www.google.com"
	password := "1234"

	repositoryMock := &repositoryMock{}
	repositoryMock.On("Save", ctx, mock.MatchedBy(func(l link.Link) bool {
		return l.URL == url && l.Password != nil
	})).Return(1)

	service := link.NewService(repositoryMock)

	// When
	l, err := service.Create(ctx, url, password)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, 1, l.ID)
}

func TestService_Redirect(t *testing.T) {
	// Given
	ctx := context.Background()
	url := "https://www.google.com"
	password := "1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	l := link.Link{
		ID:       1,
		URL:      url,
		Password: hash,
		Count:    0,
	}

	repositoryMock := &repositoryMock{}
	repositoryMock.On("FindByID", ctx, l.ID).Return(l, nil)
	repositoryMock.On("Update", ctx, mock.MatchedBy(func(l2 link.Link) bool {
		l.Count++
		return assert.ObjectsAreEqual(l, l2)
	})).Return(nil)

	service := link.NewService(repositoryMock)

	// When
	l, err = service.Redirect(ctx, 1, password)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, 1, l.ID)
	require.Equal(t, 1, l.Count)
}

func TestService_FindByID(t *testing.T) {
	// Given
	ctx := context.Background()
	l := link.Link{ID: 1}

	repositoryMock := &repositoryMock{}
	repositoryMock.On("FindByID", ctx, l.ID).Return(l, nil)

	service := link.NewService(repositoryMock)

	// When
	l, err := service.FindByID(ctx, l.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	require.Equal(t, 1, l.ID)
}

func TestService_Inactivate(t *testing.T) {
	// Given
	ctx := context.Background()
	l := link.Link{ID: 1}

	repositoryMock := &repositoryMock{}
	repositoryMock.On("FindByID", ctx, l.ID).Return(l, nil)
	repositoryMock.On("Update", ctx, mock.MatchedBy(func(l2 link.Link) bool {
		l.Inactive = true
		return assert.ObjectsAreEqual(l, l2)
	})).Return(nil)

	service := link.NewService(repositoryMock)

	// When
	err := service.Inactivate(ctx, l.ID)
	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, repositoryMock)
}
