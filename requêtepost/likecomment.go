package requêtepost

import (
	"SportHub-Forum/internal/database"
	"database/sql"
)

type LikeComment struct {
	UserID    int
	CommentID int
	Type      string
}

func AddOrUpdateLikeToComment(userID, commentID int, likeType string) error {
	query := `
		INSERT INTO like_comment (user_id, comment_id, type)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE type = VALUES(type)
	`
	_, err := database.GetDB().Exec(query, userID, commentID, likeType)
	return err
}

func RemoveLikeFromComment(userID, commentID int) error {
	query := `DELETE FROM like_comment WHERE user_id = ? AND comment_id = ?`
	_, err := database.GetDB().Exec(query, userID, commentID)
	return err
}

func GetUserLikeTypeOnComment(userID, commentID int) (string, error) {
	var likeType string
	query := `SELECT type FROM like_comment WHERE user_id = ? AND comment_id = ?`
	err := database.GetDB().QueryRow(query, userID, commentID).Scan(&likeType)
	if err == sql.ErrNoRows {
		return "", nil // Aucun like
	}
	return likeType, err
}

func CountCommentLikesByType(commentID int, likeType string) (int, error) {
	query := `SELECT COUNT(*) FROM like_comment WHERE comment_id = ? AND type = ?`
	var count int
	err := database.GetDB().QueryRow(query, commentID, likeType).Scan(&count)
	return count, err
}
