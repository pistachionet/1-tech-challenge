package services

import (
	"context"
	"database/sql/driver"
	"log/slog"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/navid/blog/internal/models"
)

func TestUsersService_ReadUser(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint64
		expectedOutput models.User
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "john", "john@me.com", "password123!"),
			mockError: nil,
			input:     1,
			expectedOutput: models.User{
				ID:       1,
				Name:     "john",
				Email:    "john@me.com",
				Password: "password123!",
			},
			expectedError: nil,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			logger := slog.Default()

			if tc.mockCalled {
				mock.
					ExpectQuery(regexp.QuoteMeta(`
                        SELECT id,
                               name,
                               email,
                               password
                        FROM users
                        WHERE id = $1::int
                    `)).
					WithArgs(tc.mockInputArgs...).
					WillReturnRows(tc.mockOutput).
					WillReturnError(tc.mockError)
			}

			userService := NewUsersService(logger, db)

			output, err := userService.ReadUser(context.TODO(), tc.input)
			if err != tc.expectedError {
				t.Errorf("expected no error, got %v", err)
			}
			if output != tc.expectedOutput {
				t.Errorf("expected %v, got %v", tc.expectedOutput, output)
			}

			if tc.mockCalled {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
