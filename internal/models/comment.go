package models

import (
	"context"
	"strings"
)

// Comment represents a comment in the system.
type Comment struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
	BlogID  uint   `json:"blog_id"`
}

// Valid checks the Comment object and returns any problems.
func (c Comment) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if strings.TrimSpace(c.Content) == "" {
		problems["content"] = "content is required"
	}

	return problems
}
