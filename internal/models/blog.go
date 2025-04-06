package models

import (
	"context"
	"strings"
	"time"
)

// Blog represents a blog in the system.
type Blog struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Score     float32   `json:"score"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_date"`
}

// Valid checks the Blog object and returns any problems.
func (b Blog) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if strings.TrimSpace(b.Title) == "" {
		problems["title"] = "title is required"
	} else if len(b.Title) > 255 {
		problems["title"] = "title cannot exceed 255 characters"
	}

	if b.UserID == 0 {
		problems["user_id"] = "user_id is required"
	}

	return problems
}
