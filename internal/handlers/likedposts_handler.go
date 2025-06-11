package handlers

import (
	"SportHub-Forum/internal/database"
	"html/template"
	"net/http"
)

// LikedPostsHandler handles displaying posts liked by the connected user
func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve posts liked by the user
	likedPosts, err := database.GetLikedPostsByUserID(database.GetDB(), userID)
	if err != nil {
		http.Error(w, "Error retrieving liked posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type PostWithLikes struct {
		ID           int
		Title        string
		Content      string
		UserID       int
		Username     string
		CategoryID   int
		CategoryName string
		CreatedAt    interface{}
		LikeCount    int
	}

	var postsWithLikes []PostWithLikes
	for _, post := range likedPosts {
		likeCount, err := database.GetLikeCountForPost(database.GetDB(), post.ID)
		if err != nil {
			http.Error(w, "Error retrieving like count: "+err.Error(), http.StatusInternalServerError)
			return
		}

		postsWithLikes = append(postsWithLikes, PostWithLikes{
			ID:           post.ID,
			Title:        post.Title,
			Content:      post.Content,
			UserID:       post.UserID,
			Username:     post.Username,
			CategoryID:   post.CategoryID,
			CategoryName: post.CategoryName,
			CreatedAt:    post.CreatedAt,
			LikeCount:    likeCount,
		})
	}

	data := struct {
		IsAuthenticated bool
		UserID          int
		LikedPosts      []PostWithLikes
	}{
		IsAuthenticated: true,
		UserID:          userID,
		LikedPosts:      postsWithLikes,
	}

	tmpl, err := template.ParseFiles("web/templates/likedposts.gohtml")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error displaying template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
