package handlers

import (
	"SportHub-Forum/internal/session"
	"net/http"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session.EndSession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
