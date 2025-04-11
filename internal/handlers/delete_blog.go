package handlers

import (
    "log/slog"
    "net/http"
    "strconv"
    "strings"

    "github.com/navid/blog/internal/services"
)

// HandleDeleteBlog handles the deletion of a blog by its ID.
func HandleDeleteBlog(logger *slog.Logger, blogsService *services.BlogService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Extract the blog ID from the URL path
        pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
        if len(pathParts) < 3 || pathParts[2] == "" {
            http.Error(w, "Blog ID not provided", http.StatusBadRequest)
            return
        }

        id, err := strconv.Atoi(pathParts[2])
        if err != nil {
            logger.ErrorContext(ctx, "failed to parse id",
                slog.String("id", pathParts[2]),
                slog.String("error", err.Error()))
            http.Error(w, "Invalid Blog ID", http.StatusBadRequest)
            return
        }

        // Delete the blog (and associated comments, if implemented)
        err = blogsService.DeleteBlog(ctx, uint(id))
        if err != nil {
            logger.ErrorContext(ctx, "failed to delete blog",
                slog.Int("id", id),
                slog.String("error", err.Error()))

            if err.Error() == "no blog found" {
                http.Error(w, "Blog not found", http.StatusNotFound)
                return
            }

            http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
            return
        }

        // Respond with success
        w.WriteHeader(http.StatusNoContent)
    }
}