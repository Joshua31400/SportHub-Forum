package handlers

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	// For POST requests (notification deletion)
	if r.Method == "POST" {
		// Check authentication
		userID, isValid := database.ValidateSession(r)
		if !isValid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
			return
		}

		action := r.FormValue("action")

		if action == "delete_all" {
			// Delete all user notifications
			err := database.DeleteAllNotificationsByUserID(userID)
			if err != nil {
				http.Error(w, "Error deleting notifications: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else if action == "delete_one" {
			// Delete a specific notification
			notificationIDStr := r.FormValue("notification_id")
			if notificationIDStr == "" {
				http.Error(w, "Missing notification ID", http.StatusBadRequest)
				return
			}

			notificationID, err := strconv.Atoi(notificationIDStr)
			if err != nil {
				http.Error(w, "Invalid notification ID", http.StatusBadRequest)
				return
			}

			// Verify that the notification belongs to the user
			notifications, err := database.GetNotificationsByUserID(userID)
			if err != nil {
				http.Error(w, "Error verifying notifications: "+err.Error(), http.StatusInternalServerError)
				return
			}

			isOwner := false
			for _, notif := range notifications {
				if notif.ID == notificationID {
					isOwner = true
					break
				}
			}

			if !isOwner {
				http.Error(w, "You are not authorized to delete this notification", http.StatusForbidden)
				return
			}

			// Delete the notification
			err = database.DeleteNotification(notificationID)
			if err != nil {
				http.Error(w, "Error deleting notification: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Unrecognized action", http.StatusBadRequest)
			return
		}

		// Redirect to notifications page after deletion
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
		return
	}

	// For GET requests (displaying notifications)
	userID, isValid := database.ValidateSession(r)
	isAuthenticated := isValid

	var notifications []models.Notification
	if isAuthenticated {
		// Retrieve user notifications
		var err error
		notifications, err = database.GetNotificationsByUserID(userID)
		if err != nil {
			http.Error(w, "Error retrieving notifications: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Prepare data for the template
	data := struct {
		IsAuthenticated bool
		Notifications   []models.Notification
	}{
		IsAuthenticated: isAuthenticated,
		Notifications:   notifications,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("web/templates/notifications.gohtml")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
