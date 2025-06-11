package handlers

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"html/template"
	"net/http"
	"strings"
)

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Verify the URL and split it for check if the post ID is valid
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Post ID wrong", http.StatusBadRequest)
		return
	}
	postID := parts[2]

	cookie, err := r.Cookie("user_id")
	isAuthenticated := err == nil

	var userID string
	if isAuthenticated {
		userID = cookie.Value
	}
	post, err := database.GetPostByID(database.GetDB(), postID)
	if err != nil {
		http.Error(w, "Error to get the post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		UserID          string
		Post            models.Post
	}{
		IsAuthenticated: isAuthenticated,
		UserID:          userID,
		Post:            post,
	}

	tmpl, err := template.ParseFiles("web/templates/post.gohtml")
	if err != nil {
		http.Error(w, "Error to charge template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error to display the template: "+err.Error(), http.StatusInternalServerError)
	}
}
