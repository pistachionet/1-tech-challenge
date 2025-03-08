package routes

import (
	"github.com/navid/blog/internal/handlers"
	"github.com/navid/blog/internal/services"
	"github.com/swaggo/http-swagger" // http-swagger middleware
	"log/slog"
	"net/http"
)

//	@title						Blog Service API
//	@version					1.0
//	@description				Practice Go API using the Standard Library and Postgres
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8000
//	@BasePath					/api
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func AddRoutes(mux *http.ServeMux, logger *slog.Logger, usersService *services.UsersService, baseURL string) {
	// Create the adapter
	userLister := handlers.NewUserListerAdapter(usersService)

	// Read a user
	mux.Handle("GET /api/users/{id}", handlers.HandleReadUser(logger, usersService))

	// Create a user
	mux.Handle("POST /api/users", handlers.HandleCreateUser(logger, usersService))

	// Update a user
	mux.Handle("PUT /api/users/{id}", handlers.HandleUpdateUser(logger, usersService))

	// Delete a user
	mux.Handle("DELETE /api/users/{id}", handlers.HandleDeleteUser(logger, usersService))

	// List users
	mux.Handle("GET /api/users", handlers.HandleListUsers(logger, userLister))

	// swagger docs
	mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(httpSwagger.URL(baseURL+"/swagger/doc.json")),
	)
	logger.Info("Swagger running", slog.String("url", baseURL+"/swagger/index.html"))
}
