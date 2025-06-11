package requêtepost

import (
	"SportHub-Forum/internal/database"
	"database/sql"
	"log"
	"time"
)

type LikePost struct {
	Log       *log.Logger
	UserID    int
	PostID    int
	Type      string
	CreatedAt time.Time
}

func AddOrUpdateLikeToPost(userID, postID int, likeType string) error {
	query := `
        INSERT INTO like_post (user_id, post_id, type)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE type = VALUES(type)
    `
	_, err := database.GetDB().Exec(query, userID, postID, likeType)
	return err
}

func RemoveLikeFromPost(userID, postID int) error {
	query := `DELETE FROM like_post WHERE user_id = ? AND post_id = ?`
	_, err := database.GetDB().Exec(query, userID, postID)
	return err
}

func GetUserLikeType(userID, postID int) (string, error) {
	var likeType string
	query := `SELECT type FROM like_post WHERE user_id = ? AND post_id = ?`
	err := database.GetDB().QueryRow(query, userID, postID).Scan(&likeType)
	if err == sql.ErrNoRows {
		return "", nil // Aucun like trouvé
	}
	return likeType, err
}

func CountLikesByType(postID int, likeType string) (int, error) {
	query := `SELECT COUNT(*) FROM like_post WHERE post_id = ? AND type = ?`
	var count int
	err := database.GetDB().QueryRow(query, postID, likeType).Scan(&count)
	return count, err
}

// Lister tous les likes d’un post (optionnel)
func GetAllLikesForPost(postID int) ([]LikePost, error) {
	query := `SELECT user_id, post_id, type FROM like_post WHERE post_id = ?`
	rows, err := database.GetDB().Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []LikePost
	for rows.Next() {
		var l LikePost
		if err := rows.Scan(&l.UserID, &l.PostID, &l.Type); err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}
	return likes, nil
}
