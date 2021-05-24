package link

import (
	"context"
	"errors"
)

type InMemoryRepository struct {
	m map[int]Link
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[int]Link),
	}
}

func (r *InMemoryRepository) Update(ctx context.Context, l Link) error {
	if l.ID == 0 {
		return errors.New("invalid link ID")
	}

	r.m[l.ID] = l
	return nil
}

func (r *InMemoryRepository) Save(ctx context.Context, l Link) int {
	l.ID = len(r.m) + 1
	r.m[l.ID] = l
	return l.ID
}

func (r *InMemoryRepository) FindByID(ctx context.Context, ID int) (Link, error) {
	link, ok := r.m[ID]
	if !ok {
		return Link{}, ErrNotFound
	}

	return link, nil
}
