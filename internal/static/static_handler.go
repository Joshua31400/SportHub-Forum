package static

import (
	"net/http"
	"path/filepath"
	"runtime"
)

func SetupStaticFiles(mux *http.ServeMux) {
	// Get the project root directory
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")

	// Static content setup
	staticDir := filepath.Join(projectRoot, "web/static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}
