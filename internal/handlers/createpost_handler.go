package handlers

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func CreatepostepageHandler(w http.ResponseWriter, r *http.Request) {
	// Verify that the request method is POST for form submission
	if r.Method == "POST" {
		r.ParseForm()

		// Get user ID from session
		userID, isValid := database.ValidateSession(r)
		if !isValid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryIDStr := r.FormValue("category")

		if title == "" || content == "" || categoryIDStr == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Convert categoryID from string to int
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		post := &models.Post{
			Title:      title,
			Content:    content,
			CategoryID: categoryID,
			UserID:     userID,
			CreatedAt:  time.Now(),
		}

		// Save the post to the database
		db := database.GetDB()
		postID, err := database.CreatePost(db, post)
		if err != nil {
			http.Error(w, "Error creating post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the newly created post page (into prinpal page)
		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
		return
	}

	_, err := r.Cookie("user_id")
	isAuthenticated := err == nil

	// Get categories from database
	repo := database.NewCategoryRepository()
	categories, err := repo.GetAll()
	if err != nil {
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := struct {
		IsAuthenticated bool
		Categories      []models.Category
	}{
		IsAuthenticated: isAuthenticated,
		Categories:      categories,
	}

	// Parse and execute template
	tmpl, err := template.ParseFiles("web/templates/createpost.gohtml")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
