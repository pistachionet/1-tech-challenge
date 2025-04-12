package routes

import (
	"log/slog"
	"net/http"

	"github.com/navid/blog/internal/handlers"
	"github.com/navid/blog/internal/services"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// @title						Blog Service API
// @version					1.0
// @description				Practice Go API using the Standard Library and Postgres
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.url				http://www.swagger.io/support
// @contact.email				support@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8000
// @BasePath					/api
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func AddRoutes(mux *http.ServeMux, logger *slog.Logger, usersService *services.UsersService, blogsService *services.BlogService, commentsService *services.CommentsService, baseURL string) {
	// User endpoints
	mux.Handle("POST /api/user", handlers.HandleCreateUser(logger, usersService))
	mux.Handle("GET /api/user", handlers.HandleListUsers(logger, handlers.NewUserListerAdapter(usersService)))
	mux.Handle("GET /api/user/{id}", handlers.HandleReadUser(logger, usersService))
	mux.Handle("PUT /api/user/{id}", handlers.HandleUpdateUser(logger, usersService))
	mux.Handle("DELETE /api/user/{id}", handlers.HandleDeleteUser(logger, usersService))

	// Blog endpoints
	logger.Info("Registering route", slog.String("method", "GET"), slog.String("path", "/api/blog"))
	mux.Handle("GET /api/blog", handlers.HandleListBlogs(logger, handlers.NewBlogListerAdapter(blogsService)))

	logger.Info("Registering route", slog.String("method", "GET"), slog.String("path", "/api/blog/{id}"))
	mux.Handle("GET /api/blog/{id}", handlers.HandleGetBlog(logger, blogsService))

	logger.Info("Registering route", slog.String("method", "PUT"), slog.String("path", "/api/blog/{id}"))
	mux.Handle("PUT /api/blog/{id}", handlers.HandleUpdateBlog(logger, blogsService, usersService))

	logger.Info("Registering route", slog.String("method", "POST"), slog.String("path", "/api/blog"))
	mux.Handle("POST /api/blog", handlers.HandleCreateBlog(logger, blogsService, usersService))

	logger.Info("Registering route", slog.String("method", "DELETE"), slog.String("path", "/api/blog/{id}"))
	mux.Handle("DELETE /api/blog/{id}", handlers.HandleDeleteBlog(logger, blogsService))

	logger.Info("Registering route", slog.String("method", "GET"), slog.String("path", "/api/comments"))
	mux.Handle("GET /api/comments", handlers.HandleListComments(logger, commentsService))

	logger.Info("Registering route", slog.String("method", "PUT"), slog.String("path", "/api/comments"))
	mux.Handle("PUT /api/comments", handlers.HandleUpdateComment(logger, commentsService, usersService, blogsService))

	logger.Info("Registering route", slog.String("method", "POST"), slog.String("path", "/api/comments"))
	mux.Handle("POST /api/comments", handlers.HandleCreateComment(logger, commentsService, usersService, blogsService))

	logger.Info("Registering route", slog.String("method", "DELETE"), slog.String("path", "/api/comments"))
	mux.Handle("DELETE /api/comments", handlers.HandleDeleteComment(logger, commentsService))

	// For debugging purposes, let's add a catch-all handler to help identify mismatched routes
	mux.HandleFunc("GET /api/blog/", func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "Caught by catch-all handler",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method))
		http.Error(w, "Route not found. Please use /api/blog/{id} format", http.StatusNotFound)
	})

	// Swagger docs
	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(httpSwagger.URL(baseURL+"/swagger/doc.json")),
	)
	logger.Info("Swagger running", slog.String("url", baseURL+"/swagger/index.html"))

	// Health check
	mux.Handle("/api/health", handlers.HandleHealthCheck(logger))
}
