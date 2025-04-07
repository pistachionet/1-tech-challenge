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

// blogReader represents a type capable of reading a blog from storage
type blogReader interface {
	GetBlog(ctx context.Context, id uint) (models.Blog, error)
}

// @Summary      Get Blog
// @Description  Get a blog by ID
// @Tags         blog
// @Produce      json
// @Param        id   path        string  true    "Blog ID"
// @Success      200  {object}    models.Blog
// @Failure      400  {object}    string
// @Failure      404  {object}    string
// @Failure      500  {object}    string
// @Router       /api/blog/{id} [get]
func HandleGetBlog(logger *slog.Logger, blogReader blogReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "Handling GET blog by ID request",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method))

		// Extract ID from path using a more compatible approach
		pathParts := strings.Split(r.URL.Path, "/")
		idStr := ""
		if len(pathParts) > 0 {
			idStr = pathParts[len(pathParts)-1]
		}
		logger.InfoContext(r.Context(), "Extracted ID from path", slog.String("id", idStr))

		if idStr == "" {
			logger.ErrorContext(r.Context(), "blog id not provided")
			http.Error(w, "Blog ID not provided", http.StatusNotFound)
			return
		}

		id64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to parse id",
				slog.String("id", idStr),
				slog.String("error", err.Error()))
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		blog, err := blogReader.GetBlog(r.Context(), uint(id64))
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to get blog",
				slog.String("error", err.Error()))

			if strings.Contains(err.Error(), "no blog found") {
				http.Error(w, "Blog not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			logger.ErrorContext(r.Context(), "failed to encode response",
				slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
