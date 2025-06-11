package handlers

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
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
	var userIDInt int
	if isAuthenticated {
		userID = cookie.Value
		userIDInt, _ = strconv.Atoi(userID)
	}

	post, err := database.GetPostByID(database.GetDB(), postID)
	if err != nil {
		http.Error(w, "Error to get the post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := database.GetCommentsByPostID(database.GetDB(), postID)
	if err != nil {
		http.Error(w, "Error to get comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	likeCount, err := database.GetLikesCountByPostID(database.GetDB(), post.ID)
	if err != nil {
		http.Error(w, "Error to get likes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify if the user has liked the post
	isLiked := false
	if isAuthenticated {
		isLiked, err = database.IsPostLikedByUser(database.GetDB(), post.ID, userIDInt)
		if err != nil {
			http.Error(w, "Error checking like status: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data := struct {
		IsAuthenticated bool
		UserID          string
		Post            models.Post
		Comments        []models.Comment
		LikeCount       int
		IsLiked         bool
	}{
		IsAuthenticated: isAuthenticated,
		UserID:          userID,
		Post:            post,
		Comments:        comments,
		LikeCount:       likeCount,
		IsLiked:         isLiked,
	}

	tmpl, err := template.ParseFiles("web/templates/post.gohtml")
	if err != nil {
		http.Error(w, "Error to charge template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error to display the template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
