package routes

import (
	"SportHub-Forum/internal/handlers"
	"SportHub-Forum/internal/middleware"
	"SportHub-Forum/internal/static"
	"net/http"
)

// Contains the setup for the HTTP routes of the application.
func SetupRoutes(mux *http.ServeMux) http.Handler {
	// Static files setup
	static.SetupStaticFiles(mux)

	// User routes
	mux.HandleFunc("/createpost", handlers.CreatepostepageHandler)
	mux.HandleFunc("/createuser", handlers.CreateUserHandler)
	mux.HandleFunc("/post/", handlers.PostPageHandler)
	mux.HandleFunc("/addcomment", handlers.AddCommentHandler)
	mux.HandleFunc("/like-post", handlers.LikePostHandler)
	mux.HandleFunc("/liked-posts", handlers.LikedPostsHandler)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/", handlers.PrincipalPageHandler)
	mux.HandleFunc("/profile", handlers.ProfilePageHandler)
	mux.HandleFunc("/filter", handlers.FilterPostsHandler)
	// Logout returns the user to the login page
	mux.HandleFunc("/logout", handlers.HandleLogout)

	return middleware.AuthMiddleware(mux)
}
