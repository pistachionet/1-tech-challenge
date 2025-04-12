package handlers

import (
    "log/slog"
    "net/http"
    "strconv"

    "github.com/navid/blog/internal/services"
)

// HandleDeleteComment handles the deletion of a comment by author_id and blog_id.
func HandleDeleteComment(logger *slog.Logger, commentsService *services.CommentsService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Parse query parameters
        authorIDStr := r.URL.Query().Get("author_id")
        blogIDStr := r.URL.Query().Get("blog_id")

        if authorIDStr == "" || blogIDStr == "" {
            http.Error(w, "author_id and blog_id are required query parameters", http.StatusBadRequest)
            return
        }

        authorID, err := strconv.Atoi(authorIDStr)
        if err != nil {
            http.Error(w, "Invalid author_id", http.StatusBadRequest)
            return
        }

        blogID, err := strconv.Atoi(blogIDStr)
        if err != nil {
            http.Error(w, "Invalid blog_id", http.StatusBadRequest)
            return
        }

        // Delete the comment
        err = commentsService.DeleteComment(ctx, authorID, blogID)
        if err != nil {
            logger.ErrorContext(ctx, "failed to delete comment", slog.String("error", err.Error()))
            if err.Error() == "no comment found" {
                http.Error(w, "Comment not found", http.StatusNotFound)
                return
            }
            http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
            return
        }

        // Respond with success
        w.WriteHeader(http.StatusNoContent)
    }
}