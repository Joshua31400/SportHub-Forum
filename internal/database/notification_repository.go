package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"time"
)

func CreateNotification(db *sql.DB, notification *models.Notification) (int64, error) {
	query := `INSERT INTO notification (userID, message, type, createdAt, postID, isRead) VALUES ( ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, notification.UserID, notification.Message, notification.IsRead, notification.PostID, notification.Type, time.Now())

	if err != nil {
		return 0, err

	}
	return result.LastInsertId()
}

// Cette fonction récupère toutes les notifications d'un utilisateur spécifique, triées par date de création (la plus récente d'abord).
func GetNotificationsByUser(db *sql.DB, userID int) ([]models.Notification, error) {
	query := `SELECT id, userID, message, type, isRead, postID, createdAt 
              FROM notification 
              WHERE userID = ? 
              ORDER BY createdAt DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.Type,
			&n.IsRead, &n.PostID, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// Dans database/notification_repository.go

// Récupère les notifications d'un utilisateur
func GetNotificationsByUser(db *sql.DB, userID int) ([]models.Notification, error) {
	query := `SELECT id, userID, message, type, createdAt, postID, isRead 
              FROM notification WHERE userID = ? ORDER BY createdAt DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.Type, &n.CreatedAt, &n.PostID, &n.IsRead); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// Marque une notification comme lue
func MarkNotificationAsRead(db *sql.DB, notificationID int) error {
	query := `UPDATE notification SET isRead = true WHERE id = ?`
	_, err := db.Exec(query, notificationID)
	return err
}

// Compte les notifications non lues d'un utilisateur
func CountUnreadNotifications(db *sql.DB, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM notification WHERE userID = ? AND isRead = false`

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}
