package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
)

type CategoryRepository struct {
	db *sql.DB
}

// Create a new instance of CategoryRepository with a database connection
func NewCategoryRepository() CategoryRepository {
	return CategoryRepository{
		db: GetDB(),
	}
}

// Get all categories from the database
func (r CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name FROM category"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving categories: %v", err)
	}
	defer rows.Close()

	categories := []models.Category{} // Init slice to hold categories
	for rows.Next() {
		var cat models.Category
		// Scan the row into the category struct
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("Error scanning categories: %v", err)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
