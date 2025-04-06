package routes

import (
	"log/slog"
	"net/http"

	"github.com/navid/blog/internal/handlers"
	"github.com/navid/blog/internal/services"
	"github.com/swaggo/http-swagger" // http-swagger middleware
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
func AddRoutes(mux *http.ServeMux, logger *slog.Logger, usersService *services.UsersService, blogsService *services.BlogService, baseURL string) {
	// User endpoints
	mux.Handle("POST /api/user", handlers.HandleCreateUser(logger, usersService))
	mux.Handle("GET /api/user", handlers.HandleListUsers(logger, handlers.NewUserListerAdapter(usersService)))
	mux.Handle("GET /api/user/{id}", handlers.HandleReadUser(logger, usersService))
	mux.Handle("PUT /api/user/{id}", handlers.HandleUpdateUser(logger, usersService))
	mux.Handle("DELETE /api/user/{id}", handlers.HandleDeleteUser(logger, usersService))

	logger.Info("Registering route", slog.String("method", "GET"), slog.String("path", "/api/blog"))
	// Blog endpoints
	mux.Handle("GET /api/blog", handlers.HandleListBlogs(logger, handlers.NewBlogListerAdapter(blogsService)))

	// Swagger docs
	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(httpSwagger.URL(baseURL+"/swagger/doc.json")),
	)
	logger.Info("Swagger running", slog.String("url", baseURL+"/swagger/index.html"))

	// Health check
	mux.Handle("/api/health", handlers.HandleHealthCheck(logger))
}
