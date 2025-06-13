package routes

import (
	"SportHub-Forum/internal/handlers"
	"SportHub-Forum/internal/middleware"
	"SportHub-Forum/internal/static"
	"net/http"
	"strings"
)

// Contains the setup for the HTTP routes of the application.
func SetupRoutes(mux *http.ServeMux) http.Handler {
	// Static files setup
	static.SetupStaticFiles(mux)

	// Initialize GitHub auth handler
	githubAuthHandler := handlers.NewGitHubAuthHandler()

	// Public routes (no authentication required)
	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/login", handlers.HandleLogin)
	publicMux.HandleFunc("/createuser", handlers.CreateUserHandler)
	publicMux.HandleFunc("/logout", handlers.HandleLogout)

	// GitHub OAuth routes
	publicMux.HandleFunc("/auth/github/login", githubAuthHandler.GitHubLogin)
	publicMux.HandleFunc("/auth/github/callback", githubAuthHandler.GitHubCallback)

	// Protected routes (authentication required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/createpost", handlers.CreatepostepageHandler)
	protectedMux.HandleFunc("/post/", handlers.PostPageHandler)
	protectedMux.HandleFunc("/addcomment", handlers.AddCommentHandler)
	protectedMux.HandleFunc("/like-post", handlers.LikePostHandler)
	protectedMux.HandleFunc("/liked-posts", handlers.LikedPostsHandler)
	protectedMux.HandleFunc("/", handlers.PrincipalPageHandler)
	protectedMux.HandleFunc("/profile", handlers.ProfilePageHandler)

	// Route dispatcher
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve static files directly
		if strings.HasPrefix(r.URL.Path, "/static/") {
			mux.ServeHTTP(w, r)
			return
		}

		// Route to public endpoints
		if isPublicRoute(r.URL.Path) {
			publicMux.ServeHTTP(w, r)
			return
		}

		// Route to protected endpoints with authentication
		middleware.AuthMiddleware(protectedMux).ServeHTTP(w, r)
	})
}

// isPublicRoute checks if a route is public (no authentication required)
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/login",
		"/createuser",
		"/logout",
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
mux.HandleFunc("/createpost", handlers.CreatepostepageHandler)
mux.HandleFunc("/createuser", handlers.CreateUserHandler)
mux.HandleFunc("/notification", handlers.NotificationHandler)
mux.HandleFunc("/post/", handlers.PostPageHandler)
mux.HandleFunc("/addcomment", handlers.AddCommentHandler)
mux.HandleFunc("/like-post", handlers.LikePostHandler)
mux.HandleFunc("/liked-posts", handlers.LikedPostsHandler)
mux.HandleFunc("/login", handlers.HandleLogin)
mux.HandleFunc("/", handlers.PrincipalPageHandler)
mux.HandleFunc("/profile", handlers.ProfilePageHandler)