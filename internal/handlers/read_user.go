package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/navid/blog/internal/models"
)

// readUserResponse represents the response for reading a user.
type readUserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userReader represents a type capable of reading a user from storage and
// returning it or an error.
type userReader interface {
	ReadUser(ctx context.Context, id uint64) (models.User, error)
}

//	@Summary		Read User
//	@Description	Read User by ID
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	readUserResponse
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/users/{id}  [GET]

func HandleReadUser(logger *slog.Logger, userReader userReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "HandleReadUser called", slog.String("path", r.URL.Path))

		ctx := r.Context()

		// Extract the "id" from the URL path
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 3 || pathParts[2] == "" {
			http.Error(w, "User ID not provided", http.StatusNotFound)
			return
		}
		idStr := pathParts[2]

		// Convert the ID from string to uint64
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			logger.ErrorContext(ctx, "failed to parse id from url", slog.String("id", idStr), slog.String("error", err.Error()))
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Read the user
		user, err := userReader.ReadUser(ctx, id)
		if err != nil {
			logger.ErrorContext(ctx, "failed to read user", slog.String("error", err.Error()))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Write the response as JSON
		response := readUserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.ErrorContext(ctx, "failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
