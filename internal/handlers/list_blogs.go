package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/navid/blog/internal/models"
	"github.com/navid/blog/internal/services"
)

// blogLister represents a type capable of listing blogs from storage and
// returning them or an error.
type blogLister interface {
	ListBlogs(ctx context.Context, title string) ([]models.Blog, error)
}

type blogListerAdapter struct {
	service *services.BlogService
}

func (a *blogListerAdapter) ListBlogs(ctx context.Context, title string) ([]models.Blog, error) {
	// Delegate to the actual service method
	return a.service.ListBlogsWithFilter(ctx, title)
}

func NewBlogListerAdapter(service *services.BlogService) blogLister {
	return &blogListerAdapter{service: service}
}

// @Summary		List Blogs
// @Description	List all blogs or filter by title
// @Tags			blog
// @Accept			json
// @Produce		json
// @Param			title	query		string	false	"Filter by title"
// @Success		200	{array}		models.Blog
// @Failure		500	{object}	string
// @Router			/blogs [GET]
func HandleListBlogs(logger *slog.Logger, blogLister blogLister) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "HandleListBlogs called", slog.String("path", r.URL.Path))

		// Get the "title" query parameter
		title := r.URL.Query().Get("title")

		// Retrieve blogs from the blogLister
		blogs, err := blogLister.ListBlogs(r.Context(), title)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to list blogs", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Write the response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(blogs); err != nil {
			logger.ErrorContext(r.Context(), "failed to encode response", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
