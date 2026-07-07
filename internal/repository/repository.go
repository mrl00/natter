// Package repository defines the Repository interface for data persistence.
// Implementations of this interface can be found under internal/infra.
package repository

import "github.com/mrl00/natter/internal/model"

// Repository defines the contract for persisting spaces and messages.
type Repository interface {
	CreateSpace(name, owner string) *model.Space
	GetSpace(id string) (*model.Space, bool)
	AddMessage(spaceID, author, content string) (*model.Message, bool)
	ListMessages(spaceID string) ([]*model.Message, bool)
	GetMessage(spaceID, messageID string) (*model.Message, bool)
}
