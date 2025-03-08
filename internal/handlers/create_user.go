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

//	@Summary		Create User
//	@Description	Create a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/users [POST]
func HandleCreateUser(logger *slog.Logger, userCreator userCreator) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            logger.ErrorContext(r.Context(), "failed to decode request body", slog.String("error", err.Error()))
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        createdUser, err := userCreator.CreateUser(r.Context(), user)
        if err != nil {
            logger.ErrorContext(r.Context(), "failed to create user", slog.String("error", err.Error()))
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(createdUser); err != nil {
            logger.ErrorContext(r.Context(), "failed to encode response", slog.String("error", err.Error()))
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    })
}