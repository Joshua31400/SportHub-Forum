package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type FilterRepository struct {
	db *sql.DB
}

// NewFilterRepository crée une nouvelle instance de FilterRepository
func NewFilterRepository(db *sql.DB) *FilterRepository {
	log.Println("Initialisation du FilterRepository")
	if db == nil {
		log.Println("ERREUR: La connexion à la base de données est nil")
	}
	return &FilterRepository{db: db}
}

// GetFilteredPosts récupère les posts selon les critères de filtrage
func (fr *FilterRepository) GetFilteredPosts(filter models.Filter) ([]models.Post, error) {
	log.Println("Début de GetFilteredPosts")

	query := `
        SELECT p.id, p.userID, p.title, p.content, p.imageURL, p.createdAt,
               COALESCE(c.id, 0), COALESCE(c.name, ''), u.userName, COUNT(lp.id) AS like_count
        FROM posts p
        LEFT JOIN like_post lp ON lp.postID = p.id
        LEFT JOIN post_categories pc ON p.id = pc.postID
        LEFT JOIN categories c ON pc.categoryID = c.id
        LEFT JOIN users u ON p.userID = u.userID
    `
	var args []interface{}
	var where string
	var orderBy string = "p.createdAt DESC"

	if filter.CategoryID != nil && *filter.CategoryID != 0 {
		where = " WHERE pc.categoryID = ? "
		args = append(args, *filter.CategoryID)
		log.Printf("Filtre par catégorie: %d", *filter.CategoryID)
	}

	if filter.SortBy != nil && *filter.SortBy == "likes" {
		orderBy = "like_count DESC"
		log.Printf("Tri par: %s", *filter.SortBy)
	}

	query += where + " GROUP BY p.id ORDER BY " + orderBy

	if filter.Limit != nil && *filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, *filter.Limit)
		log.Printf("Limite: %d", *filter.Limit)
	} else {
		query += " LIMIT 50"
		log.Println("Limite par défaut: 50")
	}

	log.Printf("Requête SQL: %s", query)
	log.Printf("Arguments: %v", args)

	rows, err := fr.db.Query(query, args...)
	if err != nil {
		errMsg := fmt.Sprintf("ERREUR lors de l'exécution de la requête: %v", err)
		log.Println(errMsg)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	log.Println("Lecture des résultats de la requête...")

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImageURL,
			&post.CreatedAt, &post.CategoryID, &post.CategoryName, &post.Username, &post.LikeCount,
		); err != nil {
			errMsg := fmt.Sprintf("ERREUR lors du scan des données: %v", err)
			log.Println(errMsg)
			return nil, err
		}
		posts = append(posts, post)
	}

	log.Printf("%d posts récupérés", len(posts))

	// Récupérer les catégories multiples pour chaque post
	log.Println("Récupération des catégories pour chaque post...")
	for i := range posts {
		log.Printf("Récupération des catégories pour le post ID: %d", posts[i].ID)
		categories, err := fr.getPostCategories(posts[i].ID)
		if err == nil {
			posts[i].Category = categories
			log.Printf("  - Post %d: %d catégories trouvées", posts[i].ID, len(categories))
		} else {
			log.Printf("  - ERREUR lors de la récupération des catégories pour le post %d: %v", posts[i].ID, err)
		}
	}

	log.Println("Fin de GetFilteredPosts")
	return posts, nil
}

// getPostCategories récupère les catégories d'un post spécifique
func (fr *FilterRepository) getPostCategories(postID int) ([]string, error) {
	log.Printf("Récupération des catégories pour le post ID: %d", postID)

	query := `
        SELECT c.name
        FROM categories c
        JOIN post_categories pc ON c.id = pc.categoryID
        WHERE pc.postID = ?
    `
	log.Printf("Requête SQL: %s avec postID=%d", query, postID)

	rows, err := fr.db.Query(query, postID)
	if err != nil {
		errMsg := fmt.Sprintf("ERREUR lors de l'exécution de la requête getPostCategories: %v", err)
		log.Println(errMsg)
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			errMsg := fmt.Sprintf("ERREUR lors du scan des catégories: %v", err)
			log.Println(errMsg)
			return nil, err
		}
		categories = append(categories, category)
		log.Printf("  - Catégorie trouvée: %s", category)
	}

	log.Printf("Total de %d catégories trouvées pour le post %d", len(categories), postID)
	return categories, nil
}

// GetAllCategories récupère toutes les catégories disponibles
func (fr *FilterRepository) GetAllCategories() ([]models.Category, error) {
	log.Println("Récupération de toutes les catégories")

	query := "SELECT id, name FROM categories"
	log.Printf("Requête SQL: %s", query)

	rows, err := fr.db.Query(query)
	if err != nil {
		errMsg := fmt.Sprintf("ERREUR lors de l'exécution de la requête GetAllCategories: %v", err)
		log.Println(errMsg)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			errMsg := fmt.Sprintf("ERREUR lors du scan des catégories: %v", err)
			log.Println(errMsg)
			return nil, err
		}
		categories = append(categories, category)
		log.Printf("  - Catégorie trouvée: ID=%d, Nom=%s", category.ID, category.Name)
	}

	log.Printf("Total de %d catégories récupérées", len(categories))
	return categories, nil
}
