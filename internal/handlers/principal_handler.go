package handlers

import (
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
)

// PrincipalPageHandler serves the main forum page
func PrincipalPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, err := r.Cookie("user_id")
	isAuthenticated := err == nil

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	templatePath := filepath.Join(projectRoot, "web/templates/principal.gohtml")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
	}{
		IsAuthenticated: isAuthenticated,
	}

	tmpl.Execute(w, data)
}
