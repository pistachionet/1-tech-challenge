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
    return &BlogService{
        db:     db,
        logger: logger,
    }
}

// CreateBlog inserts a new blog into the database.
func (s *BlogService) CreateBlog(ctx context.Context, blog models.Blog) (models.Blog, error) {
    s.logger.DebugContext(ctx, "Creating blog", "title", blog.Title)

    var createdBlog models.Blog
    err := s.db.QueryRowContext(
        ctx,
        `INSERT INTO blogs (title, score, author_id, created_date)
         VALUES ($1, $2, $3, $4)
         RETURNING id, title, score, author_id, created_date`,
        blog.Title, blog.Score, blog.UserID, blog.CreatedAt,
    ).Scan(&createdBlog.ID, &createdBlog.Title, &createdBlog.Score, &createdBlog.UserID, &createdBlog.CreatedAt)

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
        `SELECT id, title, score, author_id, created_date
         FROM blogs
         WHERE id = $1`,
        id,
    ).Scan(&blog.ID, &blog.Title, &blog.Score, &blog.UserID, &blog.CreatedAt)

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
         SET title = $1, score = $2, created_date = $3
         WHERE id = $4
         RETURNING id, title, score, author_id, created_date`,
        blog.Title, blog.Score, blog.CreatedAt, id,
    ).Scan(&updatedBlog.ID, &updatedBlog.Title, &updatedBlog.Score, &updatedBlog.UserID, &updatedBlog.CreatedAt)

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

// ListBlogsWithFilter retrieves all blogs, optionally filtering by title.
func (s *BlogService) ListBlogsWithFilter(ctx context.Context, title string) ([]models.Blog, error) {
    s.logger.DebugContext(ctx, "Listing blogs", slog.String("title", title))

    query := `SELECT id, title, score, author_id, created_date FROM blogs`
    var args []interface{}

    if title != "" {
        query += ` WHERE title ILIKE $1`
        args = append(args, "%"+title+"%")
    }

    rows, err := s.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to list blogs: %w", err)
    }
    defer rows.Close()

    var blogs []models.Blog
    for rows.Next() {
        var blog models.Blog
        if err := rows.Scan(&blog.ID, &blog.Title, &blog.Score, &blog.UserID, &blog.CreatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan blog: %w", err)
        }
        blogs = append(blogs, blog)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }

    return blogs, nil
}