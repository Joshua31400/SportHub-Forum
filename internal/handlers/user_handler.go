package handlers

import (
	"SportHub-Forum/internal/database"
	"net/http"
)

// CreateUserHandler handles user registration (GET displays form, POST processes registration)
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Manage the post request method
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Erreur d'analyse du formulaire", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are obligatory", http.StatusBadRequest)
			return
		}

		// Database creation of the user
		if err := database.CreateUser(username, email, password); err != nil {
			http.Error(w, "Failed to created the user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "web/templates/createuser.gohtml")
}
