package link_test

import (
	"context"
	"testing"

	"github.com/emacampolo/link-tracker/internal/link"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository_Update_NotFound(t *testing.T) {
	// Given
	ctx := context.Background()
	repository := link.NewInMemoryRepository()

	// When
	err := repository.Update(ctx, link.Link{ID: 1})

	// Then
	require.ErrorIs(t, err, link.ErrNotFound)
}

func TestInMemoryRepository_Update(t *testing.T) {
	// Given
	ctx := context.Background()
	repository := link.NewInMemoryRepository()
	l := newLink()
	id, err := repository.Save(ctx, l)
	if err != nil {
		t.Fatal(err)
	}

	// When
	err = repository.Update(ctx, link.Link{ID: id, Inactive: true, Count: 20})

	// Then
	require.NoError(t, err)
	l, err = repository.FindByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, id, l.ID)
	require.Equal(t, true, l.Inactive)
	require.Equal(t, 20, l.Count)
}

func TestInMemoryRepository_Save(t *testing.T) {
	// Given
	ctx := context.Background()
	repository := link.NewInMemoryRepository()
	l := newLink()

	// When
	id, err := repository.Save(ctx, l)
	if err != nil {
		t.Fatal(err)
	}

	l.ID = id

	// Then
	requireEqualLink(t, l)
}

func TestInMemoryRepository_FindByID(t *testing.T) {
	// Given
	ctx := context.Background()
	repository := link.NewInMemoryRepository()
	id, err := repository.Save(ctx, newLink())
	if err != nil {
		t.Fatal(err)
	}

	// When
	l, err := repository.FindByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	requireEqualLink(t, l)
}

func newLink() link.Link {
	return link.Link{
		URL:      "https://www.google.com",
		Password: []byte(`password`),
		Count:    10,
	}
}

func requireEqualLink(t *testing.T, l link.Link) {
	expected := newLink()
	require.Equal(t, expected.URL, l.URL)
	require.Equal(t, expected.Password, l.Password)
	require.Equal(t, expected.Count, l.Count)
	require.Equal(t, expected.Inactive, l.Inactive)
}
