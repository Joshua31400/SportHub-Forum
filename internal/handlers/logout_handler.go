package handlers

import (
	"SportHub-Forum/internal/database"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	database.EndSession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
