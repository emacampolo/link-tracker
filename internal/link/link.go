package link

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned when a Link is not found by any of its attributes.
var ErrNotFound = errors.New("link not found")

// ErrAuthentication is returned when the provided credentials cannot be validated.
var ErrAuthentication = errors.New("authentication failed")

// ErrInactive is returned when trying to redirect to a link that has been inactivated.
var ErrInactive = errors.New("link is inactive")

// Link represents an underlying URL with statistics on how it is used.
type Link struct {
	ID       int
	URL      string
	Password []byte
	Count    int
	Inactive bool
}

// Service encapsulates the business logic of a Link.
// As stated by this principle https://golang.org/doc/effective_go#generality,
// since the underlying concrete implementation does not export any other method that is not in the interface,
// we decided to define it where it is implemented rather where it is used (commonly in a handler).
type Service interface {
	Create(ctx context.Context, url, password string) (Link, error)
	Redirect(ctx context.Context, ID int, password string) (Link, error)
	FindByID(ctx context.Context, ID int) (Link, error)
	Inactivate(ctx context.Context, ID int) error
}

// Repository encapsulates the storage of a Link.
type Repository interface {
	Save(ctx context.Context, l Link) (int, error)
	Update(ctx context.Context, l Link) error
	FindByID(ctx context.Context, ID int) (Link, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(ctx context.Context, url, password string) (Link, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Link{}, err
	}

	l := Link{
		Password: hash,
		URL:      url,
	}

	id, err := s.repository.Save(ctx, l)
	if err != nil {
		return Link{}, err
	}

	l.ID = id
	return l, nil
}

func (s *service) Redirect(ctx context.Context, ID int, password string) (Link, error) {
	link, err := s.repository.FindByID(ctx, ID)
	if err != nil {
		return Link{}, ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword(link.Password, []byte(password)); err != nil {
		return Link{}, ErrAuthentication
	}

	if link.Inactive {
		return Link{}, ErrInactive
	}

	link.Count++
	if err := s.repository.Update(ctx, link); err != nil {
		return Link{}, err
	}

	return link, nil
}

func (s *service) FindByID(ctx context.Context, ID int) (Link, error) {
	return s.repository.FindByID(ctx, ID)
}

func (s *service) Inactivate(ctx context.Context, ID int) error {
	link, err := s.repository.FindByID(ctx, ID)
	if err != nil {
		return ErrNotFound
	}

	link.Inactive = true
	return s.repository.Update(ctx, link)
}
