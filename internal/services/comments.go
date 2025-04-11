package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/navid/blog/internal/models"
)

type CommentsService struct {
    db     *sql.DB
    logger *slog.Logger
}

// NewCommentsService creates a new CommentsService.
func NewCommentsService(db *sql.DB, logger *slog.Logger) *CommentsService {
    return &CommentsService{
        db:     db,
        logger: logger,
    }
}

// ListComments retrieves all comments, optionally filtering by author_id or blog_id.
func (s *CommentsService) ListComments(ctx context.Context, authorID, blogID *int) ([]models.Comment, error) {
    s.logger.DebugContext(ctx, "Listing comments", slog.Any("author_id", authorID), slog.Any("blog_id", blogID))

    query := `SELECT user_id, blog_id, message, created_date FROM comments`
    var args []interface{}
    var conditions []string

    if authorID != nil {
        conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
        args = append(args, *authorID)
    }
    if blogID != nil {
        conditions = append(conditions, fmt.Sprintf("blog_id = $%d", len(args)+1))
        args = append(args, *blogID)
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    // Log the constructed query and arguments
    s.logger.DebugContext(ctx, "Constructed query", slog.String("query", query), slog.Any("args", args))

    rows, err := s.db.QueryContext(ctx, query, args...)
    if err != nil {
        s.logger.ErrorContext(ctx, "Failed to execute query", slog.String("error", err.Error()))
        return nil, fmt.Errorf("failed to list comments: %w", err)
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var comment models.Comment
        if err := rows.Scan(&comment.UserID, &comment.BlogID, &comment.Message, &comment.CreatedDate); err != nil {
            return nil, fmt.Errorf("failed to scan comment: %w", err)
        }
        comments = append(comments, comment)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }

    return comments, nil
}