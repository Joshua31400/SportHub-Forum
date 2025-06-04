package handlers

import (
	"SportHub-Forum/internal/database"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// HandleLogin processes user login requests
// GET: serves login form, POST: authenticates user
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate required fields
		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		log.Printf("Login attempt for email: %s", email)

		// Get user from database
		user, err := database.GetUserByEmail(email)
		if err != nil {
			log.Printf("Error retrieving user: %v", err)

			// Check if user doesn't exist
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Verify password
		if user.Password != password {
			log.Printf("Invalid password for user: %s", email)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Create authentication cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    strconv.Itoa(int(user.UserID)),
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Serve login form for GET requests
	tmpl := template.Must(template.ParseFiles("web/templates/login.gohtml"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error rendering login template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
