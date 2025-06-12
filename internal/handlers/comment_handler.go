package handlers

import (
	"SportHub-Forum/internal/database"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthod refused", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "You need to be connected to comment", http.StatusUnauthorized)
		return
	}

	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")

	// Convert postIDStr to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Post ID wrong", http.StatusBadRequest)
		return
	}

	err = database.AddComment(database.GetDB(), content, postID, userID, time.Now())
	if err != nil {
		http.Error(w, "Error to add comment", http.StatusInternalServerError)
		return
	}

	if err := database.CreateCommentNotification(database.GetDB(), postID, userID, content); err != nil {
		fmt.Printf("Error to create comment notif: %v\n", err)
	}

	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}
