package routes

import (
	"SportHub-Forum/internal/handlers"
	"SportHub-Forum/internal/static"
	"net/http"
)

// Contains the setup for the HTTP routes of the application.
func SetupRoutes(mux *http.ServeMux) {
	// Static files setup
	static.SetupStaticFiles(mux)

	// User routes
	mux.HandleFunc("/createuser", handlers.CreateUserHandler)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/", handlers.PrincipalPageHandler)
	mux.HandleFunc("/profile", handlers.ProfilePageHandler)
}
