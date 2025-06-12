package handlers

import (
	"SportHub-Forum/internal/database"
	"html/template"
	"net/http"
	"strconv"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("action")

		switch action {
		case "delete":
			notificationID, err := strconv.Atoi(r.FormValue("notification_id"))
			if err != nil {
				http.Error(w, "Invalid notification ID", http.StatusBadRequest)
				return
			}
			if err := database.DeleteNotification(database.GetDB(), notificationID); err != nil {
				http.Error(w, "Error deleting notification: "+err.Error(), http.StatusInternalServerError)
				return
			}

		case "delete_all":
			if err := database.DeleteAllNotifications(database.GetDB(), userID); err != nil {
				http.Error(w, "Error deleting all notifications: "+err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Unrecognized action", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		notifications, err := database.GetNotificationsByUserID(database.GetDB(), userID)
		if err != nil {
			http.Error(w, "Error retrieving notifications: "+err.Error(), http.StatusInternalServerError)
			return
		}

		notifCount, err := database.CountNotifications(database.GetDB(), userID)
		if err != nil {
			http.Error(w, "Error counting notifications: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			IsAuthenticated bool
			UserID          int
			Notifications   interface{}
			NotifCount      int
		}{
			IsAuthenticated: true,
			UserID:          userID,
			Notifications:   notifications,
			NotifCount:      notifCount,
		}

		tmpl, err := template.ParseFiles("web/templates/notifications.gohtml")
		if err != nil {
			http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
