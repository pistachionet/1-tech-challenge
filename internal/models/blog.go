package models

import (
    "context"
    "strings"
)

// Blog represents a blog in the system.
type Blog struct {
    ID      uint   `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    UserID  uint   `json:"user_id"`
}

// Valid checks the Blog object and returns any problems.
func (b Blog) Valid(ctx context.Context) map[string]string {
    problems := make(map[string]string)

    if strings.TrimSpace(b.Title) == "" {
        problems["title"] = "title is required"
    }
    if strings.TrimSpace(b.Content) == "" {
        problems["content"] = "content is required"
    }

    return problems
}