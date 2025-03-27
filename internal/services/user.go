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
	s.logger.DebugContext(ctx, "Creating user", "user", user)

	row := s.db.QueryRowContext(
		ctx,
		`
		INSERT INTO users (name, email, password)
		VALUES ($1::int, $2::text, $3::text, $4::text)
		RETURNING id
		`,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	err := row.Scan(&user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf(
			"[in services.UsersService.CreateUser] failed to create user: %w",
			err,
		)
	}

	return user, nil
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

	result, err := s.db.ExecContext(
		ctx,
		`
	UPDATE users
	SET name = $1::text, email = $2::text, password = $3::text
	WHERE id = $4::int
	`,
		patch.Name,
		patch.Email,
		patch.Password,
		id,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("[in services.UsersServices.UpdateUser] failed to update user: %w", err)
	}

	// Checks if the user was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.User{}, fmt.Errorf("[in services.UsersService.UpdateUser] failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return models.User{}, fmt.Errorf("[in services.UsersService.UpdateUser] no user found with id: %d", id)
	}

	updatedUser := models.User{
		ID:       uint(id),
		Name:     patch.Name,
		Email:    patch.Email,
		Password: patch.Password,
	}

	return updatedUser, nil
}

// DeleteUser attempts to delete the user with the provided id. An error is
// returned if the delete fails.
func (s *UsersService) DeleteUser(ctx context.Context, id uint64) error {
	s.logger.DebugContext(ctx, "Deleting user", "id", id)

	result, err := s.db.ExecContext(
		ctx,
		`
			DELETE FROM users
			WHERE id = $1::int
			`,
		id,
	)
	if err != nil {
		return fmt.Errorf("[in services.UsersServices.DeleteUser] failed to delete user: %w", err)
	}

	// Checks if the user was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[in services.UsersService.UpdateUser] failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("[in services.UsersService.UpdateUser] no user found with id: %d", id)
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
