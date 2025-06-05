// internal/middleware/auth.go (extrait modifié)
package middleware

import (
	"SportHub-Forum/internal/session"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/createuser" ||
			r.URL.Path == "/static" || strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		userID, valid := session.ValidateSession(r)

		if !valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
