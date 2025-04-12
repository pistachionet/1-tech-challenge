package services

import (
    "context"
    "database/sql/driver"
    "fmt"
    "log/slog"
    "regexp"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/navid/blog/internal/models"
)

func TestCommentsService(t *testing.T) {
    logger := slog.Default()

    t.Run("CreateComment", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        if err != nil {
            t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
        }
        defer db.Close()

        commentsService := NewCommentsService(db, logger)

        testCases := map[string]struct {
            input          models.Comment
            mockQuery      string
            mockArgs       []driver.Value
            mockRows       *sqlmock.Rows
            mockError      error
            expectedOutput models.Comment
            expectedError  error
        }{
            "happy path": {
                input: models.Comment{
                    UserID: 1, BlogID: 2, Message: "Test Comment",
                },
                mockQuery: `INSERT INTO comments (user_id, blog_id, message, created_date) VALUES ($1, $2, $3, $4) RETURNING user_id, blog_id, message, created_date`,
                mockArgs:  []driver.Value{1, 2, "Test Comment", sqlmock.AnyArg()},
                mockRows: sqlmock.NewRows([]string{"user_id", "blog_id", "message", "created_date"}).
                    AddRow(1, 2, "Test Comment", time.Now()),
                mockError: nil,
                expectedOutput: models.Comment{
                    UserID: 1, BlogID: 2, Message: "Test Comment",
                },
                expectedError: nil,
            },
            "database error": {
                input: models.Comment{
                    UserID: 1, BlogID: 2, Message: "Test Comment",
                },
                mockQuery:      `INSERT INTO comments (user_id, blog_id, message, created_date) VALUES ($1, $2, $3, $4) RETURNING user_id, blog_id, message, created_date`,
                mockArgs:       []driver.Value{1, 2, "Test Comment", sqlmock.AnyArg()},
                mockRows:       nil,
                mockError:      fmt.Errorf("database error"),
                expectedOutput: models.Comment{},
                expectedError:  fmt.Errorf("failed to create comment: database error"),
            },
        }

        for name, tc := range testCases {
            t.Run(name, func(t *testing.T) {
                if tc.mockError != nil {
                    mock.ExpectQuery(regexp.QuoteMeta(tc.mockQuery)).
                        WithArgs(tc.mockArgs...).
                        WillReturnError(tc.mockError)
                } else {
                    mock.ExpectQuery(regexp.QuoteMeta(tc.mockQuery)).
                        WithArgs(tc.mockArgs...).
                        WillReturnRows(tc.mockRows)
                }

                output, err := commentsService.CreateComment(context.TODO(), tc.input)
                if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
                    t.Errorf("expected error %v, got %v", tc.expectedError, err)
                } else if err == nil && tc.expectedError != nil {
                    t.Errorf("expected error %v, got nil", tc.expectedError)
                } else if err != nil && tc.expectedError == nil {
                    t.Errorf("expected no error, got %v", err)
                }

                // Compare output fields except CreatedDate
                if output.UserID != tc.expectedOutput.UserID || output.BlogID != tc.expectedOutput.BlogID || output.Message != tc.expectedOutput.Message {
                    t.Errorf("expected output %v, got %v", tc.expectedOutput, output)
                }
            })
        }
    })

    t.Run("DoesCommentExist", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        if err != nil {
            t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
        }
        defer db.Close()

        commentsService := NewCommentsService(db, logger)

        testCases := map[string]struct {
            userID         int
            blogID         int
            mockQuery      string
            mockArgs       []driver.Value
            mockResult     bool
            mockError      error
            expectedOutput bool
            expectedError  error
        }{
            "comment exists": {
                userID:         1,
                blogID:         2,
                mockQuery:      `SELECT EXISTS(SELECT 1 FROM comments WHERE user_id = $1 AND blog_id = $2)`,
                mockArgs:       []driver.Value{1, 2},
                mockResult:     true,
                mockError:      nil,
                expectedOutput: true,
                expectedError:  nil,
            },
            "comment does not exist": {
                userID:         1,
                blogID:         3,
                mockQuery:      `SELECT EXISTS(SELECT 1 FROM comments WHERE user_id = $1 AND blog_id = $2)`,
                mockArgs:       []driver.Value{1, 3},
                mockResult:     false,
                mockError:      nil,
                expectedOutput: false,
                expectedError:  nil,
            },
            "database error": {
                userID:         1,
                blogID:         2,
                mockQuery:      `SELECT EXISTS(SELECT 1 FROM comments WHERE user_id = $1 AND blog_id = $2)`,
                mockArgs:       []driver.Value{1, 2},
                mockResult:     false,
                mockError:      fmt.Errorf("database error"),
                expectedOutput: false,
                expectedError:  fmt.Errorf("failed to check comment existence: database error"),
            },
        }

        for name, tc := range testCases {
            t.Run(name, func(t *testing.T) {
                if tc.mockError != nil {
                    mock.ExpectQuery(regexp.QuoteMeta(tc.mockQuery)).
                        WithArgs(tc.mockArgs...).
                        WillReturnError(tc.mockError)
                } else {
                    mock.ExpectQuery(regexp.QuoteMeta(tc.mockQuery)).
                        WithArgs(tc.mockArgs...).
                        WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(tc.mockResult))
                }

                output, err := commentsService.DoesCommentExist(context.TODO(), tc.userID, tc.blogID)
                if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
                    t.Errorf("expected error %v, got %v", tc.expectedError, err)
                } else if err == nil && tc.expectedError != nil {
                    t.Errorf("expected error %v, got nil", tc.expectedError)
                } else if err != nil && tc.expectedError == nil {
                    t.Errorf("expected no error, got %v", err)
                }

                if output != tc.expectedOutput {
                    t.Errorf("expected output %v, got %v", tc.expectedOutput, output)
                }
            })
        }
    })
}