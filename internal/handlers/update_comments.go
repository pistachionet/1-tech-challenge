package handlers

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "strconv"

    "github.com/navid/blog/internal/models"
    "github.com/navid/blog/internal/services"
)

// HandleUpdateComment handles updating a comment by author_id and blog_id.
func HandleUpdateComment(logger *slog.Logger, commentsService *services.CommentsService, usersService *services.UsersService, blogsService *services.BlogService) http.HandlerFunc {
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

        // Decode and validate the comment object
        var comment models.Comment
        if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        if comment.UserID != authorID || comment.BlogID != blogID {
            http.Error(w, "author_id and blog_id in the body must match the query parameters", http.StatusBadRequest)
            return
        }

        // Validate that the user exists
        if !usersService.DoesUserExist(ctx, authorID) {
            http.Error(w, "User not found", http.StatusBadRequest)
            return
        }

        // Validate that the blog exists
        _, err = blogsService.GetBlog(ctx, uint(blogID))
        if err != nil {
            http.Error(w, "Blog not found", http.StatusBadRequest)
            return
        }

        // Update the comment
        updatedComment, err := commentsService.UpdateComment(ctx, comment)
        if err != nil {
            logger.ErrorContext(ctx, "failed to update comment", slog.String("error", err.Error()))
            http.Error(w, "Failed to update comment", http.StatusInternalServerError)
            return
        }

        // Respond with the updated comment
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updatedComment)
    }
}