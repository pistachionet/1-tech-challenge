package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/navid/blog/internal/models"
	"github.com/navid/blog/internal/services"
)

// userLister represents a type capable of listing users from storage and
// returning them or an error.
type userLister interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

type userListerAdapter struct {
	service *services.UsersService
}

func (a *userListerAdapter) ListUsers(ctx context.Context) ([]models.User, error) {
	// Delegate to the actual service method
	return a.service.ListUsers(ctx, 0) // Pass a default value or modify as needed
}

func NewUserListerAdapter(service *services.UsersService) userLister {
	return &userListerAdapter{service: service}
}

//	@Summary		List Users
//	@Description	List all users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.User
//	@Failure		500	{object}	string
//	@Router			/users [GET]
func HandleListUsers(logger *slog.Logger, userLister userLister) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := userLister.ListUsers(r.Context())
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to list users", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logger.ErrorContext(r.Context(), "failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
