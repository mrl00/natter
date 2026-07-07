// Package service implements the business logic layer of the Natter API.
// It mediates between the HTTP handlers and the repository, enforcing rules and returning typed errors.
package service

import (
	"errors"
	"time"

	"github.com/mrl00/natter/internal/model"
	"github.com/mrl00/natter/internal/repository"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid request")
)

// Service defines the business operations for the Natter API.
type Service interface {
	CreateSpace(name, owner string) (*model.Space, error)
	AddMessage(spaceID, author, content string) (*model.Message, error)
	ListMessages(spaceID string, since time.Time) ([]*model.Message, error)
	GetMessage(spaceID, messageID string) (*model.Message, error)
}

type service struct {
	repo repository.Repository
}

func New(r repository.Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateSpace(name, owner string) (*model.Space, error) {
	space := s.repo.CreateSpace(name, owner)
	return space, nil
}

func (s *service) AddMessage(spaceID, author, content string) (*model.Message, error) {
	msg, ok := s.repo.AddMessage(spaceID, author, content)
	if !ok {
		return nil, ErrNotFound
	}
	return msg, nil
}

func (s *service) ListMessages(spaceID string, since time.Time) ([]*model.Message, error) {
	msgs, ok := s.repo.ListMessages(spaceID)
	if !ok {
		return nil, ErrNotFound
	}

	if since.IsZero() {
		return msgs, nil
	}

	var filtered []*model.Message
	for _, m := range msgs {
		if m.CreatedAt.After(since) || m.CreatedAt.Equal(since) {
			filtered = append(filtered, m)
		}
	}
	return filtered, nil
}

func (s *service) GetMessage(spaceID, messageID string) (*model.Message, error) {
	msg, ok := s.repo.GetMessage(spaceID, messageID)
	if !ok {
		return nil, ErrNotFound
	}
	return msg, nil
}
