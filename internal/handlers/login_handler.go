package handlers

import (
	"SportHub-Forum/internal/database"
	"net/http"
)

// LoginHandler handles local user login (email/password)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	if _, isValid := database.ValidateSession(r); isValid {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "web/templates/login.gohtml")
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.ServeFile(w, r, "web/templates/login.gohtml")
			return
		}

		user, err := database.AuthenticateUser(email, password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`
                <script>
                    alert('Invalid credentials');
                    window.location.href = '/login';
                </script>
            `))
			return
		}

		err = database.CreateSession(w, user.UserID)
		if err != nil {
			http.Error(w, "Connection error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
