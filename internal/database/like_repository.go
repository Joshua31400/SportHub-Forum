package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// AddLike adds or removes a like on a post (toggle functionality)
func AddLike(db *sql.DB, postID, userID int) error {
	var postExists bool
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id = ?)", postID).Scan(&postExists); err != nil {
		return fmt.Errorf("error checking post existence: %w", err)
	}
	if !postExists {
		return fmt.Errorf("post with ID %d does not exist", postID)
	}

	// Check if user exists
	var userExists bool
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE userid = ?)", userID).Scan(&userExists); err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !userExists {
		return fmt.Errorf("user with ID %d does not exist", userID)
	}

	// Check if like already exists
	liked, err := IsPostLikedByUser(db, postID, userID)
	if err != nil {
		return fmt.Errorf("error checking existing like: %w", err)
	}

	// Toggle like status
	if liked {
		if _, err := db.Exec("DELETE FROM `like` WHERE postid = ? AND userid = ?", postID, userID); err != nil {
			return fmt.Errorf("error removing like: %w", err)
		}
	} else {
		if _, err := db.Exec("INSERT INTO `like` (postid, userid) VALUES (?, ?)", postID, userID); err != nil {
			return fmt.Errorf("error adding like: %w", err)
		}
	}

	return nil
}

// IsPostLikedByUser checks if a post is liked by a specific user
func IsPostLikedByUser(db *sql.DB, postID, userID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM `like` WHERE postid = ? AND userid = ?)"
	if err := db.QueryRow(query, postID, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if post is liked: %w", err)
	}
	return exists, nil
}

// GetLikesCountByPostID returns the number of likes for a specific post
func GetLikesCountByPostID(db *sql.DB, postID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM `like` WHERE postid = ?"
	if err := db.QueryRow(query, postID).Scan(&count); err != nil {
		return 0, fmt.Errorf("error retrieving like count: %w", err)
	}
	return count, nil
}

// GetLikeCountForPost is an alias for GetLikesCountByPostID to maintain compatibility
func GetLikeCountForPost(db *sql.DB, postID int) (int, error) {
	return GetLikesCountByPostID(db, postID)
}

// Retrieves all posts liked by a specific user
func GetLikedPostsByUserID(db *sql.DB, userID int) ([]models.Post, error) {
	query := `
  SELECT p.id, p.title, p.content, p.categoryid, c.name as categoryName,
         p.userid, u.userName, p.createdat
  FROM post p
  JOIN ` + "`like`" + ` l ON p.id = l.postid
  LEFT JOIN category c ON p.categoryid = c.id
  LEFT JOIN user u ON p.userid = u.userID
  WHERE l.userid = ?
  ORDER BY l.createdat DESC
 `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving liked posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		loc = time.UTC
	}

	for rows.Next() {
		var post models.Post
		var username, categoryName sql.NullString
		var createdAtStr string

		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryID,
			&categoryName, &post.UserID, &username, &createdAtStr); err != nil {
			return nil, fmt.Errorf("error scanning post data: %w", err)
		}

		if username.Valid {
			post.Username = username.String
		} else {
			post.Username = "Unknown user"
		}
		if categoryName.Valid {
			post.CategoryName = categoryName.String
		}

		formats := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339}
		for _, format := range formats {
			if parsed, err := time.Parse(format, createdAtStr); err == nil {
				post.CreatedAt = parsed.In(loc)
				break
			}
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through posts: %w", err)
	}
	return posts, nil
}
