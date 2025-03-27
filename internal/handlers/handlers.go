package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"

	"github.com/navid/blog/internal/services"
)

// validator is an object that can be validated.
type validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}

// decodeValid decodes a model from an http request and performs validation
// on it.
func decodeValid[T validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}

// HandleListUsersWithFilter handles the GET /api/user endpoint.
func HandleListUsersWithFilter(logger *slog.Logger, userLister *services.UsersService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the "name" query parameter
		name := r.URL.Query().Get("name")

		// Retrieve users from the service
		users, err := userLister.ListUsersWithFilter(r.Context(), name)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to list users", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Write the response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logger.ErrorContext(r.Context(), "failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
