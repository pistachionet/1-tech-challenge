package models

import (
    "context"
    "strings"
)

// User represents a user in the system.
type User struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Valid checks the User object and returns any problems.
func (u User) Valid(ctx context.Context) map[string]string {
    problems := make(map[string]string)

    if strings.TrimSpace(u.Name) == "" {
        problems["name"] = "name is required"
    }
    if strings.TrimSpace(u.Email) == "" {
        problems["email"] = "email is required"
    }
    if strings.TrimSpace(u.Password) == "" {
        problems["password"] = "password is required"
    }

    return problems
}