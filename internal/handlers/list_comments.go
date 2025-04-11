package handlers

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "strconv"

    "github.com/navid/blog/internal/services"
)

// HandleListComments handles retrieving all comments, optionally filtering by author_id or blog_id.
func HandleListComments(logger *slog.Logger, commentsService *services.CommentsService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Parse query parameters
        authorIDStr := r.URL.Query().Get("author_id")
        blogIDStr := r.URL.Query().Get("blog_id")

        var authorID, blogID *int
        if authorIDStr != "" {
            id, err := strconv.Atoi(authorIDStr)
            if err != nil {
                http.Error(w, "Invalid author_id", http.StatusBadRequest)
                return
            }
            authorID = &id
        }
        if blogIDStr != "" {
            id, err := strconv.Atoi(blogIDStr)
            if err != nil {
                http.Error(w, "Invalid blog_id", http.StatusBadRequest)
                return
            }
            blogID = &id
        }

        // Retrieve comments
        comments, err := commentsService.ListComments(ctx, authorID, blogID)
        if err != nil {
            logger.ErrorContext(ctx, "failed to list comments", slog.String("error", err.Error()))
            http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
            return
        }

        // Respond with comments
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(comments)
    }
}