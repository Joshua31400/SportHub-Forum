package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"time"
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

func GetAllPosts(db *sql.DB) ([]models.Post, error) {
	query := `
        SELECT p.id, p.title, p.content, p.categoryid, p.userid, p.createdat, u.userName
        FROM post p
        LEFT JOIN user u ON p.userid = u.userID
        ORDER BY p.createdat DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		loc = time.UTC
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var createdAtStr string
		var username sql.NullString

		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryID,
			&post.UserID, &createdAtStr, &username); err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}

		if username.Valid {
			post.Username = username.String
		} else {
			post.Username = "User is not registered"
		}

		if createdAtStr != "" {
			formats := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339}
			for _, format := range formats {
				if parsed, err := time.Parse(format, createdAtStr); err == nil {
					post.CreatedAt = parsed.In(loc)
					break
				}
			}
		}
		posts = append(posts, post)
	}
	return posts, nil
}
