package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/navid/blog/internal/models"
)

type BlogService struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewBlogService creates a new BlogService.
func NewBlogService(db *sql.DB, logger *slog.Logger) *BlogService {
	return &BlogService{db: db, logger: logger}
}

// CreateBlog inserts a new blog into the database.
func (s *BlogService) CreateBlog(ctx context.Context, blog models.Blog) (models.Blog, error) {
	s.logger.DebugContext(ctx, "Creating blog", "title", blog.Title)

	var createdBlog models.Blog
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO blogs (title, content, user_id, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, title, content, user_id, created_at, updated_at`,
		blog.Title, blog.Content, blog.UserID, blog.CreatedAt, blog.UpdatedAt,
	).Scan(&createdBlog.ID, &createdBlog.Title, &createdBlog.Content, &createdBlog.UserID, &createdBlog.CreatedAt, &createdBlog.UpdatedAt)

	if err != nil {
		return models.Blog{}, fmt.Errorf("failed to create blog: %w", err)
	}

	return createdBlog, nil
}

// GetBlog retrieves a blog by its ID.
func (s *BlogService) GetBlog(ctx context.Context, id uint) (models.Blog, error) {
	s.logger.DebugContext(ctx, "Retrieving blog", "id", id)

	var blog models.Blog
	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, title, content, user_id, created_at, updated_at
         FROM blogs
         WHERE id = $1`,
		id,
	).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.UserID, &blog.CreatedAt, &blog.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.Blog{}, fmt.Errorf("no blog found with id: %d", id)
	} else if err != nil {
		return models.Blog{}, fmt.Errorf("failed to retrieve blog: %w", err)
	}

	return blog, nil
}

// UpdateBlog updates an existing blog in the database.
func (s *BlogService) UpdateBlog(ctx context.Context, id uint, blog models.Blog) (models.Blog, error) {
	s.logger.DebugContext(ctx, "Updating blog", "id", id)

	var updatedBlog models.Blog
	err := s.db.QueryRowContext(
		ctx,
		`UPDATE blogs
         SET title = $1, content = $2, updated_at = $3
         WHERE id = $4
         RETURNING id, title, content, user_id, created_at, updated_at`,
		blog.Title, blog.Content, blog.UpdatedAt, id,
	).Scan(&updatedBlog.ID, &updatedBlog.Title, &updatedBlog.Content, &updatedBlog.UserID, &updatedBlog.CreatedAt, &updatedBlog.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.Blog{}, fmt.Errorf("no blog found with id: %d", id)
	} else if err != nil {
		return models.Blog{}, fmt.Errorf("failed to update blog: %w", err)
	}

	return updatedBlog, nil
}

// DeleteBlog deletes a blog by its ID.
func (s *BlogService) DeleteBlog(ctx context.Context, id uint) error {
	s.logger.DebugContext(ctx, "Deleting blog", "id", id)

	result, err := s.db.ExecContext(ctx, `DELETE FROM blogs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete blog: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no blog found with id: %d", id)
	}

	return nil
}
