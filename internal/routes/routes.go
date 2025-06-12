package routes

import (
	"SportHub-Forum/internal/handlers"
	"SportHub-Forum/internal/middleware"
	"SportHub-Forum/internal/static"
	"net/http"
	"strings"
)

// SetupRoutes configures all application routes
func SetupRoutes(mux *http.ServeMux) http.Handler {
	static.SetupStaticFiles(mux)

	googleAuthHandler := handlers.NewAuthHandler()
	githubAuthHandler := handlers.NewGitHubAuthHandler() // Nouveau

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/login", handlers.LoginHandler)
	publicMux.HandleFunc("/createuser", handlers.CreateUserHandler)
	publicMux.HandleFunc("/logout", handlers.HandleLogout)

	publicMux.HandleFunc("/auth/google/login", googleAuthHandler.GoogleLogin)
	publicMux.HandleFunc("/auth/google/callback", googleAuthHandler.GoogleCallback)

	publicMux.HandleFunc("/auth/github/login", githubAuthHandler.GitHubLogin)
	publicMux.HandleFunc("/auth/github/callback", githubAuthHandler.GitHubCallback)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/", handlers.PrincipalPageHandler)
	protectedMux.HandleFunc("/profile", handlers.ProfilePageHandler)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			mux.ServeHTTP(w, r)
			return
		}

		if isPublicRoute(r.URL.Path) {
			publicMux.ServeHTTP(w, r)
			return
		}

		middleware.AuthMiddleware(protectedMux).ServeHTTP(w, r)
	})
}

// isPublicRoute checks if a route is public (no authentication required)
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/login",
		"/createuser",
		"/logout",
		"/auth/google/login",
		"/auth/google/callback",
		"/auth/github/login",
		"/auth/github/callback",
	}

	for _, route := range publicRoutes {
		if path == route {
			return true
		}
	}
	return false
}
