package models

import (
	"context"
	"testing"
)

func TestUser_Valid(t *testing.T) {
	testcases := map[string]struct {
		input    User
		expected map[string]string
	}{
		"valid user": {
			input: User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			expected: map[string]string{},
		},
		"missing name": {
			input: User{
				ID:       2,
				Name:     "",
				Email:    "jane.doe@example.com",
				Password: "password123",
			},
			expected: map[string]string{
				"name": "name is required",
			},
		},
		"missing email": {
			input: User{
				ID:       3,
				Name:     "Jane Doe",
				Email:    "",
				Password: "password123",
			},
			expected: map[string]string{
				"email": "email is required",
			},
		},
		"missing password": {
			input: User{
				ID:       4,
				Name:     "Jane Doe",
				Email:    "jane.doe@example.com",
				Password: "",
			},
			expected: map[string]string{
				"password": "password is required",
			},
		},
		"missing all fields": {
			input: User{
				ID:       5,
				Name:     "",
				Email:    "",
				Password: "",
			},
			expected: map[string]string{
				"name":     "name is required",
				"email":    "email is required",
				"password": "password is required",
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			problems := tc.input.Valid(context.Background())
			if len(problems) != len(tc.expected) {
				t.Errorf("expected %d problems, got %d", len(tc.expected), len(problems))
			}
			for field, expectedMsg := range tc.expected {
				if msg, ok := problems[field]; !ok || msg != expectedMsg {
					t.Errorf("expected problem for field %s: %s, got: %s", field, expectedMsg, msg)
				}
			}
		})
	}
}
