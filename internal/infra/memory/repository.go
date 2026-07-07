// Package memory provides an in-memory implementation of the repository.Repository interface.
// Data is stored in maps protected by a sync.RWMutex and is lost when the process exits.
package memory

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/mrl00/natter/internal/model"
)

// Repository is a thread-safe in-memory implementation of repository.Repository.
type Repository struct {
	mu       sync.RWMutex
	spaces   map[string]*model.Space
	messages map[string][]*model.Message
}

func New() *Repository {
	return &Repository{
		spaces:   make(map[string]*model.Space),
		messages: make(map[string][]*model.Message),
	}
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (r *Repository) CreateSpace(name, owner string) *model.Space {
	space := &model.Space{
		ID:        generateID(),
		Name:      name,
		Owner:     owner,
		CreatedAt: time.Now().UTC(),
	}

	r.mu.Lock()
	r.spaces[space.ID] = space
	r.mu.Unlock()

	return space
}

func (r *Repository) GetSpace(id string) (*model.Space, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	space, ok := r.spaces[id]
	return space, ok
}

func (r *Repository) AddMessage(spaceID, author, content string) (*model.Message, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.spaces[spaceID]; !ok {
		return nil, false
	}

	msg := &model.Message{
		ID:        generateID(),
		SpaceID:   spaceID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	r.messages[spaceID] = append(r.messages[spaceID], msg)
	return msg, true
}

func (r *Repository) ListMessages(spaceID string) ([]*model.Message, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.spaces[spaceID]; !ok {
		return nil, false
	}
	return r.messages[spaceID], true
}

func (r *Repository) GetMessage(spaceID, messageID string) (*model.Message, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.spaces[spaceID]; !ok {
		return nil, false
	}

	for _, m := range r.messages[spaceID] {
		if m.ID == messageID {
			return m, true
		}
	}
	return nil, false
}
