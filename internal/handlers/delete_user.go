package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
)

// userDeleter represents a type capable of deleting a user from storage and
// returning an error if something goes wrong.
type userDeleter interface {
	DeleteUser(ctx context.Context, id uint64) error
}

// @Summary		Delete User
// @Description	Delete a user by ID
// @Tags			user
// @Param			id	path		string	true	"User ID"
// @Success		204	{object}	nil
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/users/{id} [DELETE]
func HandleDeleteUser(logger *slog.Logger, userDeleter userDeleter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to parse id from url", slog.String("id", idStr), slog.String("error", err.Error()))
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if err := userDeleter.DeleteUser(r.Context(), id); err != nil {
			logger.ErrorContext(r.Context(), "failed to delete user", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
