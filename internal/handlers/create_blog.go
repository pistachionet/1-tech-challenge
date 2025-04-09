package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/navid/blog/internal/models"
	"github.com/navid/blog/internal/services"
)

func HandleCreateBlog(logger *slog.Logger, blogsService *services.BlogService, usersService *services.UsersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blog models.Blog
		if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate the blog object
		if blog.Title == "" || blog.Score <= 0 {
			http.Error(w, "Invalid blog data: title and score are required", http.StatusBadRequest)
			return
		}

		// Validate author_id
		if !usersService.DoesUserExist(r.Context(), blog.AuthorID) {
			http.Error(w, "Invalid author_id: user does not exist", http.StatusBadRequest)
			return
		}

		// Create the blog
		createdBlog, err := blogsService.CreateBlog(r.Context(), blog)
		if err != nil {
			logger.Error("Failed to create blog", slog.String("error", err.Error()))
			http.Error(w, "Failed to create blog", http.StatusInternalServerError)
			return
		}

		// Respond with the created blog
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBlog)
	}
}
