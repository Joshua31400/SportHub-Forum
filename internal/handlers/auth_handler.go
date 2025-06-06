package handlers

import (
	"SportHub-Forum/internal/authentification"
	"SportHub-Forum/internal/database"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Get user from database
		user, err := database.GetUserByEmail(email)
		if err != nil {
			log.Printf("Error retrieving user: %v", err)
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the password matches
		if !authentification.CheckPasswordHash(password, user.Password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = database.CreateSession(w, user.UserID)
		if err != nil {
			http.Error(w, "Error for session creation", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get request method
	tmpl := template.Must(template.ParseFiles("web/templates/login.gohtml"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
