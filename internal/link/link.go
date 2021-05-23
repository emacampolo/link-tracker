package link

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("link not found")
var ErrAuthentication = errors.New("authentication failed")

type Link struct {
	ID       int
	URL      string
	Password []byte
	Count    int
}

// Service encapsulates the business logic and persistence of a Link.
// The implementation is not thread safe.
type Service struct {
	m map[int]Link
}

func NewService() *Service {
	return &Service{
		m: make(map[int]Link),
	}
}

func (s *Service) Create(url, password string) (Link, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Link{}, err
	}

	l := Link{
		ID:       len(s.m) + 1,
		Password: hash,
		URL:      url,
	}

	s.m[l.ID] = l
	return l, nil
}

func (s *Service) Redirect(ID int, password string) (Link, error) {
	link, ok := s.m[ID]
	if !ok {
		return Link{}, ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword(link.Password, []byte(password)); err != nil {
		return Link{}, ErrAuthentication
	}

	link.Count++
	return link, nil
}
