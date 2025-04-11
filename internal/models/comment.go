package models

import (
	"context"
	"strings"
	"time"
)

// Comment represents a comment in the system.
type Comment struct {
	UserID      int       `json:"user_id"`
	BlogID      int       `json:"blog_id"`
	Message     string    `json:"message"`
	CreatedDate time.Time `json:"created_date"`
}

// Valid checks the Comment object and returns any problems.
func (c Comment) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if strings.TrimSpace(c.Message) == "" {
		problems["Message"] = "Message is required"
	}

	return problems
}
