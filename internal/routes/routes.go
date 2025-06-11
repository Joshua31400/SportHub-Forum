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
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/", handlers.PrincipalPageHandler)
	mux.HandleFunc("/profile", handlers.ProfilePageHandler)

	// Connect to the database
	notifHandler := &handlers.NotificationHandler{DB: db}

	// Routes pour les notifications
	router.HandleFunc("/api/notifications", notifHandler.GetUserNotifications).Methods("GET")
	router.HandleFunc("/api/notifications/read", notifHandler.MarkAsRead).Methods("PUT")
	router.HandleFunc("/api/notifications/unread/count", notifHandler.GetUnreadCount).Methods("GET")

	// Logout returns the user to the login page
	mux.HandleFunc("/logout", handlers.HandleLogout)

	return middleware.AuthMiddleware(mux)
}
