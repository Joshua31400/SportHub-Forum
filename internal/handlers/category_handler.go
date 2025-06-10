package handlers

import (
	"SportHub-Forum/internal/database"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	repo database.CategoryRepository
}

// This functions get all categories from the database and returns them as JSON
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter) {
	categories, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
