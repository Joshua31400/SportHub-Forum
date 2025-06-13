package middleware

import (
	"SportHub-Forum/internal/database"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware validates user sessions and protects routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow public access to authentication routes and static files
		if r.URL.Path == "/login" || r.URL.Path == "/createuser" ||
			r.URL.Path == "/static" || strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		userID, valid := database.ValidateSession(r)

		if !valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
