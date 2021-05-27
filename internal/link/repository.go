package link

import (
	"context"
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
	_, ok := r.m[l.ID]
	if !ok {
		return ErrNotFound
	}

	r.m[l.ID] = l
	return nil
}

func (r *InMemoryRepository) Save(ctx context.Context, l Link) (int, error) {
	l.ID = len(r.m) + 1
	r.m[l.ID] = l
	return l.ID, nil
}

func (r *InMemoryRepository) FindByID(ctx context.Context, ID int) (Link, error) {
	link, ok := r.m[ID]
	if !ok {
		return Link{}, ErrNotFound
	}

	return link, nil
}
