package main

import (
	"net/http"
	"os"
	"path/filepath"
)

// setupStaticFileServer serves files from the built frontend directory (frontend/dist).
// For development this serves files from disk. It falls back to index.html for
// client-side routing when a file is not found.
func setupStaticFileServer() http.Handler {
	// compute path relative to this package (cmd/meshcomm)
	distPath := filepath.Clean(filepath.Join("..", "..", "frontend", "dist"))

	// If the dist folder doesn't exist, return a simple NotFound handler
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "frontend not built - run 'npm run build' in frontend", http.StatusNotFound)
		})
	}

	fs := http.FileServer(http.Dir(distPath))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for root and unknown routes to support SPA routing
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			return
		}

		// try to open file; if not exists, serve index.html
		file := filepath.Join(distPath, filepath.Clean(r.URL.Path))
		if _, err := os.Stat(file); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	})
}
