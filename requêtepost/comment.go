package requêtepost

import (
	"SportHub-Forum/internal/database"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
}

type CommentWithLikes struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
	UserVote  string // "like", "dislike", ou ""
}

func AddComment(postID, userID int, content string) error {
	query := `INSERT INTO comment (post_id, user_id, contenu, created_at)
	          VALUES (?, ?, ?, ?)`
	_, err := database.GetDB().Exec(query, postID, userID, content, time.Now())
	return err
}

func DeleteComment(commentID int) error {
	query := `DELETE FROM comment WHERE id = ?`
	_, err := database.GetDB().Exec(query, commentID)
	return err
}

// 🔍 Récupérer un commentaire par son ID
func GetCommentByID(commentID int) (*Comment, error) {
	query := `SELECT id, post_id, user_id, contenu, created_at FROM comment WHERE id = ?`
	row := database.GetDB().QueryRow(query, commentID)

	var c Comment
	err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetCommentsWithLikes(postID int, currentUserID int) ([]CommentWithLikes, error) {
	query := `SELECT id, post_id, user_id, contenu, created_at FROM comment WHERE post_id = ?`
	rows, err := database.GetDB().Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentWithLikes

	for rows.Next() {
		var c CommentWithLikes
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}

		c.Likes, _ = CountCommentLikesByType(c.ID, "like")
		c.Dislikes, _ = CountCommentLikesByType(c.ID, "dislike")

		c.UserVote, _ = GetUserLikeTypeOnComment(currentUserID, c.ID)

		comments = append(comments, c)
	}

	return comments, nil
}
