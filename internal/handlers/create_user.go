package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/navid/blog/internal/models"
)

// userCreator represents a type capable of creating a user in storage and
// returning it or an error.
type userCreator interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

// @Summary		Create User
// @Description	Create a new user
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			user	body		models.User	true	"User"
// @Success		201		{object}	models.User
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/users [POST]
func HandleCreateUser(logger *slog.Logger, userCreator userCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, problems, err := decodeValid[models.User](r)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to decode request body",
				slog.String("error", err.Error()))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(problems) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(problems); err != nil {
				logger.ErrorContext(r.Context(), "failed to encode validation problems",
					slog.String("error", err.Error()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		createdUser, err := userCreator.CreateUser(r.Context(), user)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to create user",
				slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Only log the important result
		logger.InfoContext(r.Context(), "user created",
			slog.Uint64("id", uint64(createdUser.ID)))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdUser); err != nil {
			logger.ErrorContext(r.Context(), "failed to encode response",
				slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
