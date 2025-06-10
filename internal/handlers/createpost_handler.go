package handlers

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"html/template"
	"net/http"
)

func CreatepostepageHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication by verifying if user_id cookie exists
	_, err := r.Cookie("user_id")
	isAuthenticated := err == nil

	// Get categories from database
	repo := database.NewCategoryRepository()
	categories, err := repo.GetAll()
	if err != nil {
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := struct {
		IsAuthenticated bool
		Categories      []models.Category
	}{
		IsAuthenticated: isAuthenticated,
		Categories:      categories,
	}

	// Parse and execute template
	tmpl, err := template.ParseFiles("web/templates/createpost.gohtml")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
