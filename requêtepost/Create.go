package requêtepost

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"log"
	"time"
)

func CreatePost(post models.Post) error {
	query := `INSERT INTO post (user_id, titre, contenu, image_url, created_at) 
	          VALUES (?, ?, ?, ?, ?)`

	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}

	_, err := database.DB.Exec(query, post.UserID, post.Titre, post.Contenu, post.ImageURL, post.CreatedAt)
	if err != nil {
		log.Println("❌ Erreur lors de l’insertion :", err)
		return err
	}

	log.Println("✅ Post inséré avec succès depuis requêtepost")
	return nil
}
