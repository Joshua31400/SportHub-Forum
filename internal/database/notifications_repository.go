package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"time"
)

func CreateLikeNotification(db *sql.DB, postID int, likerID int) error {
	var authorID int
	var postTitle string
	query := "SELECT userid, title FROM post WHERE id = ?"
	if err := db.QueryRow(query, postID).Scan(&authorID, &postTitle); err != nil {
		return fmt.Errorf("error retrieving post author: %w", err)
	}

	if authorID == likerID {
		return nil
	}

	var likerUsername string
	if err := db.QueryRow("SELECT username FROM user WHERE userid = ?", likerID).Scan(&likerUsername); err != nil {
		return fmt.Errorf("error retrieving username: %w", err)
	}

	message := fmt.Sprintf("%s a aimé votre post : %s", likerUsername, postTitle)

	return CreateNotification(db, authorID, message, "like", postID)
}

func CreateCommentNotification(db *sql.DB, postID int, commenterID int, commentContent string) error {
	var authorID int
	var postTitle string
	query := "SELECT userid, title FROM post WHERE id = ?"
	if err := db.QueryRow(query, postID).Scan(&authorID, &postTitle); err != nil {
		return fmt.Errorf("error retrieving post author: %w", err)
	}

	if authorID == commenterID {
		return nil
	}

	var commenterUsername string
	if err := db.QueryRow("SELECT username FROM user WHERE userid = ?", commenterID).Scan(&commenterUsername); err != nil {
		return fmt.Errorf("error retrieving username: %w", err)
	}

	message := fmt.Sprintf("%s a commenté votre post : %s", commenterUsername, postTitle)

	return CreateNotification(db, authorID, message, "comment", postID)
}

func CreateNotification(db *sql.DB, userID int, message string, notificationType string, postID int) error {
	var userExists bool
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE userid = ?)", userID).Scan(&userExists); err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !userExists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}
	query := "INSERT INTO notifications (userid, message, type, postid, createdat) VALUES (?, ?, ?, ?, NOW())"
	if _, err := db.Exec(query, userID, message, notificationType, postID); err != nil {
		return fmt.Errorf("error creating notification: %w", err)
	}
	return nil
}

func GetNotificationsByUserID(db *sql.DB, userID int) ([]models.Notification, error) {
	query := `
  SELECT n.id, n.userid, n.message, n.type, p.id, n.createdat
  FROM notifications n
  LEFT JOIN post p ON n.postid = p.id
  WHERE n.userid = ?
  ORDER BY n.createdat DESC
 `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		loc = time.UTC
	}

	for rows.Next() {
		var notification models.Notification
		var postID int
		var createdAtStr string

		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message,
			&notification.Type, &postID, &createdAtStr); err != nil {
			return nil, fmt.Errorf("error reading notification data: %w", err)
		}

		notification.PostID = models.Post{ID: postID}

		formats := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339}
		for _, format := range formats {
			if parsed, err := time.Parse(format, createdAtStr); err == nil {
				notification.CreatedAt = parsed.In(loc)
				break
			}
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through notifications: %w", err)
	}
	return notifications, nil
}

func DeleteNotification(db *sql.DB, notificationID int) error {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM notifications WHERE id = ?)"
	if err := db.QueryRow(checkQuery, notificationID).Scan(&exists); err != nil {
		return fmt.Errorf("error checking notification existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("no notification found with ID %d", notificationID)
	}
	query := "DELETE FROM notifications WHERE id = ?"
	_, err := db.Exec(query, notificationID)
	if err != nil {
		return fmt.Errorf("error deleting notification: %w", err)
	}
	return nil
}

func DeleteAllNotifications(db *sql.DB, userID int) error {
	query := "DELETE FROM notifications WHERE userid = ?"
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error deleting all notifications: %w", err)
	}
	return nil
}

func CountNotifications(db *sql.DB, userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM notifications WHERE userid = ?"
	if err := db.QueryRow(query, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("error counting notifications: %w", err)
	}
	return count, nil
}
