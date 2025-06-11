package handlers

import (
	"SportHub-Forum/internal/database"
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type NotificationHandler struct {
	DB *sql.DB
}

// Structure pour les données du template
type NotificationPageData struct {
	Notifications []database.Notification
	UnreadCount   int
}

// Fonction helper pour extraire l'email de l'utilisateur depuis la session
func getUserEmailFromRequest(r *http.Request) (string, error) {
	// Récupérer l'email depuis le cookie de session
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Fonction pour formater les dates dans le template
func formatDate(t time.Time) string {
	return t.Format("02/01/2006 à 15:04")
}

// Récupère toutes les notifications de l'utilisateur connecté
func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'email utilisateur depuis la session
	email, err := getUserEmailFromRequest(r)
	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	// Obtenir l'utilisateur complet
	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	notifications, err := database.GetNotificationsByUser(h.DB, user.UserID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// Marque une notification comme lue
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'email utilisateur depuis la session
	email, err := getUserEmailFromRequest(r)
	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	// Vérifier l'authentification de l'utilisateur
	_, err = database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de notification invalide", http.StatusBadRequest)
		return
	}

	err = database.MarkNotificationAsRead(h.DB, id)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour de la notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Récupère le nombre de notifications non lues
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'email utilisateur depuis la session
	email, err := getUserEmailFromRequest(r)
	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	// Obtenir l'utilisateur complet
	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	count, err := database.CountUnreadNotifications(h.DB, user.UserID)
	if err != nil {
		http.Error(w, "Erreur lors du comptage des notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

// Handler pour afficher la page des notifications
func (h *NotificationHandler) NotificationsPage(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'email utilisateur depuis la session
	email, err := getUserEmailFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtenir l'utilisateur complet
	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer les notifications
	notifications, err := database.GetNotificationsByUser(h.DB, user.UserID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des notifications", http.StatusInternalServerError)
		return
	}

	// Récupérer le nombre de notifications non lues
	count, err := database.CountUnreadNotifications(h.DB, user.UserID)
	if err != nil {
		http.Error(w, "Erreur lors du comptage des notifications", http.StatusInternalServerError)
		return
	}

	// Créer les fonctions personnalisées pour le template
	funcMap := template.FuncMap{
		"formatDate": formatDate,
	}

	// Préparer les données pour le template
	data := NotificationPageData{
		Notifications: notifications,
		UnreadCount:   count,
	}

	// Charger et exécuter le template avec les fonctions personnalisées
	tmpl, err := template.New("notifications.html").Funcs(funcMap).ParseFiles("templates/notifications.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur d'exécution du template", http.StatusInternalServerError)
	}
}

// Handler pour marquer une notification comme lue et rediriger vers la page des notifications
func (h *NotificationHandler) MarkAsReadAndRedirect(w http.ResponseWriter, r *http.Request) {
	// Appeler d'abord le handler existant pour marquer comme lu
	h.MarkAsRead(w, r)

	// Vérifier s'il n'y a pas eu d'erreur (code 200 OK)
	if w.Header().Get("Content-Type") == "" {
		// Rediriger vers la page des notifications
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
	}
	// En cas d'erreur, le handler MarkAsRead a déjà envoyé la réponse d'erreur
}
