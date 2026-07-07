// Package postgres provides a PostgreSQL implementation of the repository.Repository interface.
// Data is persisted in PostgreSQL and survives process restarts.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mrl00/natter/internal/model"
	"github.com/mrl00/natter/internal/repository"
)

var _ repository.Repository = (*Repository)(nil)

func generateID() string {
	return uuid.New().String()
}

// Repository is a PostgreSQL-backed implementation of repository.Repository.
type Repository struct {
	db *sql.DB
}

// New creates a new PostgreSQL repository and verifies the connection.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Ping verifies the database connection is alive.
func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Migrate creates the required database tables.
func (r *Repository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS spaces (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			owner      TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);
		CREATE TABLE IF NOT EXISTS messages (
			id         TEXT PRIMARY KEY,
			space_id   TEXT NOT NULL REFERENCES spaces(id),
			author     TEXT NOT NULL,
			content    TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_messages_space_id ON messages(space_id);
	`)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

func (r *Repository) CreateSpace(name, owner string) *model.Space {
	id := generateID()
	now := time.Now().UTC()

	_, _ = r.db.ExecContext(
		context.Background(),
		`INSERT INTO spaces (id, name, owner, created_at) VALUES ($1, $2, $3, $4)`,
		id, name, owner, now,
	)

	return &model.Space{
		ID:        id,
		Name:      name,
		Owner:     owner,
		CreatedAt: now,
	}
}

func (r *Repository) GetSpace(id string) (*model.Space, bool) {
	var s model.Space
	err := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, name, owner, created_at FROM spaces WHERE id = $1`, id,
	).Scan(&s.ID, &s.Name, &s.Owner, &s.CreatedAt)
	if err != nil {
		return nil, false
	}
	return &s, true
}

func (r *Repository) AddMessage(spaceID, author, content string) (*model.Message, bool) {
	var exists bool
	_ = r.db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM spaces WHERE id = $1)`, spaceID,
	).Scan(&exists)
	if !exists {
		return nil, false
	}

	id := generateID()
	now := time.Now().UTC()

	_, _ = r.db.ExecContext(
		context.Background(),
		`INSERT INTO messages (id, space_id, author, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		id, spaceID, author, content, now,
	)

	return &model.Message{
		ID:        id,
		SpaceID:   spaceID,
		Author:    author,
		Content:   content,
		CreatedAt: now,
	}, true
}

func (r *Repository) ListMessages(spaceID string) ([]*model.Message, bool) {
	var exists bool
	_ = r.db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM spaces WHERE id = $1)`, spaceID,
	).Scan(&exists)
	if !exists {
		return nil, false
	}

	rows, err := r.db.QueryContext(
		context.Background(),
		`SELECT id, space_id, author, content, created_at FROM messages WHERE space_id = $1 ORDER BY created_at`, spaceID,
	)
	if err != nil {
		return nil, false
	}

	defer func() {
		_ = rows.Close()
	}()

	var msgs []*model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.SpaceID, &m.Author, &m.Content, &m.CreatedAt); err != nil {
			return nil, false
		}
		msgs = append(msgs, &m)
	}
	return msgs, true
}

func (r *Repository) GetMessage(spaceID, messageID string) (*model.Message, bool) {
	var m model.Message
	err := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, space_id, author, content, created_at FROM messages WHERE space_id = $1 AND id = $2`,
		spaceID, messageID,
	).Scan(&m.ID, &m.SpaceID, &m.Author, &m.Content, &m.CreatedAt)
	if err != nil {
		return nil, false
	}
	return &m, true
}
