// Package model defines the core domain types shared across all layers of the Natter API.
package model

import "time"

// Space represents a social space where users can gather and exchange messages.
type Space struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"createdAt"`
}

// Message represents a single message posted within a space.
type Message struct {
	ID        string    `json:"id"`
	SpaceID   string    `json:"spaceId"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateSpaceRequest is the payload for creating a new space.
type CreateSpaceRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

// CreateMessageRequest is the payload for posting a new message to a space.
type CreateMessageRequest struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}
