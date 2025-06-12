package handlers

import (
	"SportHub-Forum/internal/database"
	"fmt"
	"net/http"
	"strconv"
)

// Handles adding or removing likes on posts
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "You must be logged in to like posts", http.StatusUnauthorized)
		return
	}

	// Get post ID from request
	var postIDStr string
	if r.Method == http.MethodPost {
		postIDStr = r.FormValue("post_id")
	} else {
		postIDStr = r.URL.Query().Get("id")
	}

	// Convert post ID to integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Toggle like status (add or remove)
	if err := database.AddLike(database.GetDB(), postID, userID); err != nil {
		http.Error(w, fmt.Sprintf("Error processing like: %v", err), http.StatusInternalServerError)
		return
	}

	if err := database.CreateLikeNotification(database.GetDB(), postID, userID); err != nil {
		fmt.Printf("Error to create like notif: %v\n", err)
	}
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}
