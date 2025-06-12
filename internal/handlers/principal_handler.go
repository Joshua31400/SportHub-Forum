package handlers

import (
	"SportHub-Forum/internal/database"
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
)

// PrincipalPageHandler serves the main forum page
func PrincipalPageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	userID, isAuthenticated := database.ValidateSession(r)

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
		UserID          int
	}{
		IsAuthenticated: isAuthenticated,
		UserID:          userID,
	}

	tmpl.Execute(w, data)
}
