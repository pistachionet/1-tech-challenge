package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/navid/blog/internal/models"
)

func TestBlogService_ReadBlog(t *testing.T) {
	testcases := map[string]struct {
		mockCalled     bool
		mockInputArgs  []driver.Value
		mockOutput     *sqlmock.Rows
		mockError      error
		input          uint
		expectedOutput models.Blog
		expectedError  error
	}{
		"happy path": {
			mockCalled:    true,
			mockInputArgs: []driver.Value{1},
			mockOutput: sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"}).
				AddRow(1, "Test Blog", "This is a test blog.", 1, parseTime("2024-05-15T10:00:00Z"), parseTime("2024-05-15T10:00:00Z")),
			mockError: nil,
			input:     1,
			expectedOutput: models.Blog{
				ID:        1,
				Title:     "Test Blog",
				Score:     1,
				UserID:    1,
				CreatedAt: parseTime("2024-05-15T10:00:00Z"),
			},
			expectedError: nil,
		},
		"blog not found": {
			mockCalled:     true,
			mockInputArgs:  []driver.Value{2},
			mockOutput:     sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"}), // No rows
			mockError:      nil,
			input:          2,
			expectedOutput: models.Blog{},
			expectedError:  fmt.Errorf("no blog found with id: %d", 2),
		},
		"database error": {
			mockCalled:     true,
			mockInputArgs:  []driver.Value{3},
			mockOutput:     nil,             // No rows should be returned
			mockError:      sql.ErrConnDone, // Simulate a database connection error
			input:          3,
			expectedOutput: models.Blog{},
			expectedError:  fmt.Errorf("failed to retrieve blog: %w", sql.ErrConnDone),
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
				query := regexp.QuoteMeta(`
                    SELECT id, title, content, user_id, created_at, updated_at
                    FROM blogs
                    WHERE id = $1
                `)
				if tc.mockError != nil {
					mock.ExpectQuery(query).
						WithArgs(tc.mockInputArgs...).
						WillReturnError(tc.mockError) // Simulate the error
				} else {
					mock.ExpectQuery(query).
						WithArgs(tc.mockInputArgs...).
						WillReturnRows(tc.mockOutput) // Return rows if no error
				}
			}

			blogService := NewBlogService(db, logger)

			output, err := blogService.GetBlog(context.TODO(), tc.input)
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
			if output != tc.expectedOutput {
				t.Errorf("expected output %v, got %v", tc.expectedOutput, output)
			}

			if tc.mockCalled {
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func parseTime(value string) time.Time {
	t, _ := time.Parse(time.RFC3339, value)
	return t
}
