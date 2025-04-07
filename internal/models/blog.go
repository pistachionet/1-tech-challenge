package models

import (
	"context"
	"strings"
	"time"
)

// Blog represents a blog in the system.
type Blog struct {
	ID        uint      `json:"id,omitempty"`
	Title     string    `json:"title" validate:"required"`
	Score     float64   `json:"score"` // Ensure this matches the database type
	UserID    uint      `json:"author_id" validate:"required"`
	CreatedAt time.Time `json:"created_date"` // Ensure this matches the database type
}

// Valid checks the Blog object and returns any problems.
func (b Blog) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if strings.TrimSpace(b.Title) == "" {
		problems["title"] = "title is required"
	}

	if b.UserID == 0 {
		problems["author_id"] = "author_id is required"
	}

	return problems
}
