package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"time"
)

func AddComment(db *sql.DB, content string, postID, userID int, createdAt time.Time) error {
	// Verify if the post exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error verifying post: %w", err)
	}
	if !exists {
		return fmt.Errorf("post with ID %d does not exist", postID)
	}

	// MySQL date format
	formattedDate := createdAt.Format("2006-01-02 15:04:05")

	// Add the comment to the database
	query := `INSERT INTO comment (content, postid, userid, createdat) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, content, postID, userID, formattedDate)
	if err != nil {
		return fmt.Errorf("error adding comment: %w", err)
	}

	return nil
}

func GetCommentsByPostID(db *sql.DB, postID string) ([]models.Comment, error) {
	query := `
  SELECT c.id, c.content, c.postid, c.userid, u.userName, c.createdat
  FROM comment c
  LEFT JOIN user u ON c.userid = u.userID
  WHERE c.postid = ?
  ORDER BY c.createdat ASC
 `

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving comments: %w", err)
	}
	defer rows.Close()

	comments := []models.Comment{}
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		loc = time.UTC
	}

	for rows.Next() {
		var comment models.Comment
		var username sql.NullString
		var createdAtStr string

		if err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID,
			&comment.UserID, &username, &createdAtStr); err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}

		if username.Valid {
			comment.Username = username.String
		} else {
			comment.Username = "Unknown user"
		}

		// Process the date
		formats := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339}
		for _, format := range formats {
			if parsed, err := time.Parse(format, createdAtStr); err == nil {
				comment.CreatedAt = parsed.In(loc)
				break
			}
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through comments: %w", err)
	}

	return comments, nil
}
