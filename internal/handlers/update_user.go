package handlers

import (
    "context"
    "encoding/json"
    "log/slog"
    "net/http"
    "strconv"

    "github.com/navid/blog/internal/models"
)

// userUpdater represents a type capable of updating a user in storage and
// returning it or an error.
type userUpdater interface {
    UpdateUser(ctx context.Context, id uint64, user models.User) (models.User, error)
}

//	@Summary		Update User
//	@Description	Update an existing user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"User ID"
//	@Param			user	body		models.User	true	"User"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	string
//	@Failure		404		{object}	string
//	@Failure		500		{object}	string
//	@Router			/users/{id} [PUT]
func HandleUpdateUser(logger *slog.Logger, userUpdater userUpdater) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        idStr := r.URL.Query().Get("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            logger.ErrorContext(r.Context(), "failed to parse id from url", slog.String("id", idStr), slog.String("error", err.Error()))
            http.Error(w, "Invalid ID", http.StatusBadRequest)
            return
        }

        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            logger.ErrorContext(r.Context(), "failed to decode request body", slog.String("error", err.Error()))
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        updatedUser, err := userUpdater.UpdateUser(r.Context(), id, user)
        if err != nil {
            logger.ErrorContext(r.Context(), "failed to update user", slog.String("error", err.Error()))
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
            logger.ErrorContext(r.Context(), "failed to encode response", slog.String("error", err.Error()))
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    })
}