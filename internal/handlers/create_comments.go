package handlers

import (
    "encoding/json"
    "log/slog"
    "net/http"

    "github.com/navid/blog/internal/models"
    "github.com/navid/blog/internal/services"
)

// HandleCreateComment handles the creation of a new comment.
func HandleCreateComment(logger *slog.Logger, commentsService *services.CommentsService, usersService *services.UsersService, blogsService *services.BlogService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Decode and validate the comment object
        var comment models.Comment
        if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Validate the comment object
        if problems := comment.Valid(ctx); len(problems) > 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(problems)
            return
        }

        // Validate that the user exists
        if !usersService.DoesUserExist(ctx, comment.UserID) {
            http.Error(w, "User not found", http.StatusBadRequest)
            return
        }

        // Validate that the blog exists
        _, err := blogsService.GetBlog(ctx, uint(comment.BlogID))
        if err != nil {
            http.Error(w, "Blog not found", http.StatusBadRequest)
            return
        }

        // Check if a comment with the same user_id and blog_id already exists
        exists, err := commentsService.DoesCommentExist(ctx, comment.UserID, comment.BlogID)
        if err != nil {
            logger.ErrorContext(ctx, "failed to check comment existence", slog.String("error", err.Error()))
            http.Error(w, "Failed to validate comment", http.StatusInternalServerError)
            return
        }
        if exists {
            http.Error(w, "Comment already exists for the given user_id and blog_id", http.StatusBadRequest)
            return
        }

        // Create the comment
        createdComment, err := commentsService.CreateComment(ctx, comment)
        if err != nil {
            logger.ErrorContext(ctx, "failed to create comment", slog.String("error", err.Error()))
            http.Error(w, "Failed to create comment", http.StatusInternalServerError)
            return
        }

        // Respond with the created comment
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdComment)
    }
}