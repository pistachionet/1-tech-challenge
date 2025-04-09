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

// blogUpdater represents a type capable of updating a blog in storage
type blogUpdater interface {
	UpdateBlog(ctx context.Context, id uint, blog models.Blog) (models.Blog, error)
}

// @Summary		Update Blog
// @Description	Update an existing blog
// @Tags			blog
// @Accept			json
// @Produce		json
// @Param			id		path		string		true	"Blog ID"
// @Param			blog	body		models.Blog	true	"Blog"
// @Success		200		{object}	models.Blog
// @Failure		400		{object}	string
// @Failure		404		{object}	string
// @Failure		500		{object}	string
// @Router			/blog/{id} [put]
func HandleUpdateBlog(logger *slog.Logger, blogUpdater blogUpdater, userReader userReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get id from path using built-in PathValue
		idStr := r.PathValue("id")
		if idStr == "" {
			http.Error(w, "Blog ID not provided", http.StatusNotFound)
			return
		}

		id64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			logger.ErrorContext(ctx, "failed to parse id",
				slog.String("id", idStr),
				slog.String("error", err.Error()))
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		id := uint(id64)

		// Decode and validate the blog
		blog, problems, err := decodeValid[models.Blog](r)
		if err != nil {
			logger.ErrorContext(ctx, "failed to decode request body",
				slog.String("error", err.Error()))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(problems) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(problems); err != nil {
				logger.ErrorContext(ctx, "failed to encode validation problems",
					slog.String("error", err.Error()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// Validate that the author exists
		_, err = userReader.ReadUser(ctx, uint64(blog.AuthorID))
		if err != nil {
			logger.ErrorContext(ctx, "author validation failed",
				slog.Uint64("author_id", uint64(blog.AuthorID)),
				slog.String("error", err.Error()))
			http.Error(w, "Author not found", http.StatusBadRequest)
			return
		}

		// Update the blog
		updatedBlog, err := blogUpdater.UpdateBlog(ctx, id, blog)
		if err != nil {
			logger.ErrorContext(ctx, "failed to update blog",
				slog.String("error", err.Error()))

			if strings.Contains(err.Error(), "no blog found") {
				http.Error(w, "Blog not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to update blog", http.StatusInternalServerError)
			return
		}

		// Return the updated blog
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(updatedBlog); err != nil {
			logger.ErrorContext(ctx, "failed to encode response",
				slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
