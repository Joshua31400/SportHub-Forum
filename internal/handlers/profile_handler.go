package handlers

import (
	"SportHub-Forum/internal/database"
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
)

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	// Recover the user ID from the context to redirect to login if not authenticated already
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userPosts, err := database.GetPostsByUserID(database.GetDB(), userID)
	if err != nil {
		http.Error(w, "Error to get posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := database.GetUserByID(database.GetDB(), userID)
	if err != nil {
		http.Error(w, "Error to get user informations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	templatePath := filepath.Join(projectRoot, "web/templates/profile.gohtml")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	// Prepare data for the template
	data := struct {
		IsAuthenticated bool
		UserID          int
		Username        string
		CreatedAt       interface{}
		UserPosts       interface{}
	}{
		IsAuthenticated: true,
		UserID:          userID,
		Username:        user.Username,
		CreatedAt:       user.CreatedAt,
		UserPosts:       userPosts,
	}
	tmpl.Execute(w, data)
}
