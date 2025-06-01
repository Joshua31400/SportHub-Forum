package handlers

import (
	"net/http"
	"path/filepath"
	"runtime"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the project template directory
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	templatePath := filepath.Join(projectRoot, "web/templates/createuser.gohtml")

	http.ServeFile(w, r, templatePath)
}
