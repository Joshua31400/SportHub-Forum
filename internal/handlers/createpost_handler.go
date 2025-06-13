package handlers

import (
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func CreatepostepageHandler(w http.ResponseWriter, r *http.Request) {
	// Verify that the request method is POST for form submission
	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20) // Limit 10MB
		if err != nil {
			http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
			return
		}

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

		var imageURL string
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			// Get current working directory
			currentDir, err := os.Getwd()
			if err != nil {
				http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Create absolute path for uploads directory
			uploadDir := filepath.Join(currentDir, "web", "static", "uploads")

			// Create directory with all parents if needed
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				http.Error(w, "Error creating uploads directory: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Verify directory exists after creation
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				http.Error(w, "Error: unable to create uploads directory", http.StatusInternalServerError)
				return
			}

			// Create unique filename
			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
			imagePath := filepath.Join(uploadDir, filename)

			// Create the file on the server
			dst, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Error creating file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the destination
			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Error saving image: "+err.Error(), http.StatusInternalServerError)
				return
			}
			// URL path for the browser (this stays relative)
			imageURL = "/static/uploads/" + filename
		}

		post := &models.Post{
			Title:      title,
			Content:    content,
			CategoryID: categoryID,
			UserID:     userID,
			CreatedAt:  time.Now(),
			ImageURL:   imageURL,
		}

		// Save the post to the database
		db := database.GetDB()
		postID, err := database.CreatePost(db, post)
		if err != nil {
			http.Error(w, "Error creating post: "+err.Error(), http.StatusInternalServerError)
			return
		}

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
