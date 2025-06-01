package handlers

import (
	"net/http"
	"path/filepath"
	"runtime"

	"SportHub-Forum/internal/database"
)

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

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))
		return
	}

	// For GET requests, serve the HTML template for user creation
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	templatePath := filepath.Join(projectRoot, "web/templates/createuser.gohtml")

	http.ServeFile(w, r, templatePath)
}
