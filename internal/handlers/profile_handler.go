package handlers

import (
	"SportHub-Forum/internal/database"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
)

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	// Recover the user ID from the context to redirect to login if not authenticated already
	userID, isAuthenticated := database.ValidateSession(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Convert userID to string for template rendering
	userIDStr := fmt.Sprintf("%d", userID)
	isAuthenticated = true

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
		UserId          string
	}{
		IsAuthenticated: isAuthenticated,
		UserId:          userIDStr,
	}
	tmpl.Execute(w, data)
}
