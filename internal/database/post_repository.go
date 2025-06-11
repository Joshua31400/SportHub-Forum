package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
)

func CreatePost(db *sql.DB, post *models.Post) (int, error) {
	result, err := db.Exec("INSERT INTO post (title, content, categoryid, userid, createdat) VALUES (?, ?, ?, ?, ?)",
		post.Title, post.Content, post.CategoryID, post.UserID, post.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("error inserting post: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting post ID: %w", err)
	}
	return int(lastID), nil
}
