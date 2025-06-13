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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postIDStr := strconv.Itoa(postID)

	userID, isAuthenticated := database.ValidateSession(r)
	db := database.GetDB()

	if r.Method == "POST" && r.FormValue("action") == "delete" {
		if !isAuthenticated {
			http.Error(w, "You must be logged in to delete this post", http.StatusUnauthorized)
			return
		}

		// Get the post to verify the author
		post, err := database.GetPostByID(db, postIDStr)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Check if the user is the author
		if userID != post.UserID {
			http.Error(w, "You are not authorized to delete this post", http.StatusForbidden)
			return
		}

		// Delete the post
		err = database.DeletePost(db, postID)
		if err != nil {
			http.Error(w, "Error deleting post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to homepage
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get post data
	post, err := database.GetPostByID(db, postIDStr)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	isAuthor := isAuthenticated && userID == post.UserID

	comments, err := database.GetCommentsByPostID(db, postIDStr)
	if err != nil {
		http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
		return
	}

	likeCount, err := database.GetLikesCountByPostID(db, postID)
	if err != nil {
		http.Error(w, "Error retrieving like count", http.StatusInternalServerError)
		return
	}

	isLiked := false
	if isAuthenticated {
		liked, err := database.IsPostLikedByUser(db, postID, userID)
		if err == nil {
			isLiked = liked
		}
	}

	data := struct {
		IsAuthenticated bool
		UserID          int
		Post            models.Post
		Comments        []models.Comment
		LikeCount       int
		IsLiked         bool
		IsAuthor        bool
		HasImage        bool
	}{
		IsAuthenticated: isAuthenticated,
		UserID:          userID,
		Post:            post,
		Comments:        comments,
		LikeCount:       likeCount,
		IsLiked:         isLiked,
		IsAuthor:        isAuthor,
		HasImage:        post.ImageURL != "",
	}

	tmpl, err := template.ParseFiles("web/templates/post.gohtml")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
