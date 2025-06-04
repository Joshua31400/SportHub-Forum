package handlers

import (
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
)

// ProfilePageHandler serves the user profile page
func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userId := cookie.Value
	isAuthenticated := true

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	templatePath := filepath.Join(projectRoot, "web/templates/profile.gohtml")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		UserId          string
	}{
		IsAuthenticated: isAuthenticated,
		UserId:          userId,
	}
	tmpl.Execute(w, data)
}
