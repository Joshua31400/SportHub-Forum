package handlers

import (
	"SportHub-Forum/internal/database"
	"net/http"
	"strconv"
	"time"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Vous devez être connecté pour commenter", http.StatusUnauthorized)
		return
	}

	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")

	// Convert postIDStr to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	err = database.AddComment(database.GetDB(), content, postID, userID, time.Now())
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}
