package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicPaths := map[string]bool{
			"/login":      true,
			"/createuser": true,
			"/static/":    true,
		}

		path := r.URL.Path
		if publicPaths[path] || (len(path) >= 8 && path[:8] == "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		_, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
