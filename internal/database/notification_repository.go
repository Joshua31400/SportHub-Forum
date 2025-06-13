package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func CreateNotification(userID int, message string, sourceType string, sourceID int) error {
	query := `INSERT INTO notification (userid, message, createdat, sourcetype, sourceid)
             VALUES (?, ?, ?, ?, ?)`

	result, err := ExecWithTimeout(query, userID, message, time.Now(), sourceType, sourceID)
	if err != nil {
		log.Printf("SQL error when creating notification: %v", err)
		return fmt.Errorf("failed to create notification: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected, notification not created")
	}
	return nil
}

func GetNotificationsByUserID(userID int) ([]models.Notification, error) {
	query := `SELECT id, userid, message, createdat, sourcetype, sourceid
             FROM notification
             WHERE userid = ?
             ORDER BY createdat DESC`

	rows, err := GetDB().Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching notifications: %v", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var createdAtStr string
		var sourceType sql.NullString
		var sourceID sql.NullInt64

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Message,
			&createdAtStr,
			&sourceType,
			&sourceID,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning notification: %v", err)
		}

		if createdAtStr != "" {
			if parsed, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
				notification.CreatedAt = parsed
			}
		}

		if sourceType.Valid {
			notification.SourceType = sourceType.String
		}

		if sourceID.Valid {
			notification.SourceID = int(sourceID.Int64)
		}

		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func DeleteNotification(notificationID int) error {
	query := `DELETE FROM notification WHERE id = ?`

	result, err := ExecWithTimeout(query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

func DeleteAllNotificationsByUserID(userID int) error {
	query := `DELETE FROM notification WHERE userid = ?`

	result, err := ExecWithTimeout(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notifications: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Deleted %d notifications for user %d", rowsAffected, userID)
	return nil
}
