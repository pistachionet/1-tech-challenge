package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/navid/blog/internal/models"
)

// UsersService is a service capable of performing CRUD operations for
// models.User models.
type UsersService struct {
	logger *slog.Logger
	db     *sql.DB
}

// NewUsersService creates a new UsersService and returns a pointer to it.
func NewUsersService(logger *slog.Logger, db *sql.DB) *UsersService {
	return &UsersService{
		logger: logger,
		db:     db,
	}
}

// CreateUser attempts to create the provided user, returning a fully hydrated
// models.User or an error.
func (s *UsersService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
    s.logger.DebugContext(ctx, "Creating user", "email", user.Email)

    var createdUser models.User
    err := s.db.QueryRowContext(
        ctx,
        `
        INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, name, email, password
        `,
        user.Name,
        user.Email,
        user.Password,
    ).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.Password)

    if err != nil {
        return models.User{}, fmt.Errorf(
            "[in services.UsersService.CreateUser] failed to create user: %w",
            err,
        )
    }

    return createdUser, nil
}

// ReadUser attempts to read a user from the database using the provided id. A
// fully hydrated models.User or error is returned.
func (s *UsersService) ReadUser(ctx context.Context, id uint64) (models.User, error) {
	s.logger.DebugContext(ctx, "Reading user", "id", id)

	row := s.db.QueryRowContext(
		ctx,
		`
		SELECT id,
		       name,
		       email,
		       password
		FROM users
		WHERE id = $1::int
        `,
		id,
	)

	var user models.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, nil
		default:
			return models.User{}, fmt.Errorf(
				"[in services.UsersService.ReadUser] failed to read user: %w",
				err,
			)
		}
	}

	return user, nil
}

// UpdateUser attempts to perform an update of the user with the provided id,
// updating, it to reflect the properties on the provided patch object. A
// models.User or an error.
func (s *UsersService) UpdateUser(ctx context.Context, id uint64, patch models.User) (models.User, error) {
    s.logger.DebugContext(ctx, "Updating user", "id", id)

    var updatedUser models.User
    err := s.db.QueryRowContext(
        ctx,
        `
        UPDATE users 
        SET name = $2, email = $3, password = $4
        WHERE id = $1
        RETURNING id, name, email, password
        `,
        id,
        patch.Name,
        patch.Email, 
        patch.Password,
    ).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Password)

    if err != nil {
        if err == sql.ErrNoRows {
            return models.User{}, fmt.Errorf("no user found with id: %d", id)
        }
        return models.User{}, fmt.Errorf("failed to update user: %w", err)
    }

    return updatedUser, nil
}

// DeleteUser attempts to delete the user with the provided id. An error is
// returned if the delete fails.
func (s *UsersService) DeleteUser(ctx context.Context, id uint64) error {
    s.logger.DebugContext(ctx, "Deleting user and related data", "id", id)

    // Start a transaction to ensure all or nothing deletion
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }
    
    // Ensure transaction either commits or rolls back
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // Delete related comments first
    _, err = tx.ExecContext(ctx, 
        `DELETE FROM comments WHERE user_id = $1`, id)
    if err != nil {
        return fmt.Errorf("failed to delete user comments: %w", err)
    }

    // Delete related blogs
    _, err = tx.ExecContext(ctx, 
        `DELETE FROM blogs WHERE user_id = $1`, id)
    if err != nil {
        return fmt.Errorf("failed to delete user blogs: %w", err)
    }

    // Finally delete the user
    result, err := tx.ExecContext(ctx, 
        `DELETE FROM users WHERE id = $1`, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    // Check if user was found
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no user found with id: %d", id)
    }

    // Commit the transaction
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

// ListUsers attempts to list all users in the database. A slice of models.User
// or an error is returned.
func (s *UsersService) ListUsers(ctx context.Context, id uint64) ([]models.User, error) {
	s.logger.DebugContext(ctx, "Listing users", "id", id)

	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT *
		FROM users
		`,
	)
	if err != nil {
		return []models.User{}, fmt.Errorf("[in services.UsersServices.ListUser] failed to list users: %w", err)
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return []models.User{}, fmt.Errorf("[in services.UsersService.ListUsers] failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return []models.User{}, fmt.Errorf("[in services.UsersService.ListUsers] rows iteration error: %w", err)
	}

	return users, nil
}

// ListUsersWithFilter retrieves all users from the database, optionally filtering by name.
func (s *UsersService) ListUsersWithFilter(ctx context.Context, name string) ([]models.User, error) {
	s.logger.DebugContext(ctx, "Listing users with filter", "name", name)

	query := `
        SELECT id, name, email, password
        FROM users
    `
	var rows *sql.Rows
	var err error

	if name != "" {
		// Add a WHERE clause to filter by name
		query += "WHERE name ILIKE $1::text"
		rows, err = s.db.QueryContext(ctx, query, "%"+name+"%")
	} else {
		// Query all users if no name filter is provided
		rows, err = s.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("[in services.UsersService.ListUsersWithFilter] failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, fmt.Errorf("[in services.UsersService.ListUsersWithFilter] failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[in services.UsersService.ListUsersWithFilter] rows iteration error: %w", err)
	}

	return users, nil
}
