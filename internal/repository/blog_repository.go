package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/navid/blog/internal/models"
)

type BlogRepository struct {
	db *sql.DB
}

// GetBlog retrieves a blog from the database by its ID
func (r *BlogRepository) GetBlog(ctx context.Context, id uint) (models.Blog, error) {
	query := `SELECT id, author_id, title, score, created_date FROM blogs WHERE id = ?`

	var blog models.Blog
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&blog.ID,
		&blog.AuthorID,
		&blog.Title,
		&blog.Score,       // Ensure this matches the database type
		&blog.CreatedDate, // Ensure this matches the database type
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Blog{}, fmt.Errorf("no blog found with id %d", id)
		}
		return models.Blog{}, fmt.Errorf("failed to retrieve blog: %w", err)
	}

	return blog, nil
}
