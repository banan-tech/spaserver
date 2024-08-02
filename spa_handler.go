package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func spaHandler(staticDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Check if the path is a directory
		fullPath := filepath.Join(staticDir, path)
		if isDirectory(fullPath) {
			dirIndexPath := filepath.Join(fullPath, "index.html")
			if fileExists(dirIndexPath) {
				// Serve index.html inside the directory if it exists
				http.ServeFile(w, r, dirIndexPath)
				return
			}
		}

		// Check if the file exists
		if fileExists(fullPath) {
			// Serve the file if it exists
			http.ServeFile(w, r, fullPath)
			return
		}

		// If the path ends with a file extension other than .html, return 404
		if ext := filepath.Ext(path); ext != "" {
			http.NotFound(w, r)
			return
		}

		// Serve the root index.html for all other cases
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	}
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Helper function to check if a path is a directory
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
