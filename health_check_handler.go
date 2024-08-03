package main

import (
	"net/http"
	"path/filepath"
)

func healthcheckHandler(rootDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dirIndexPath := filepath.Join(rootDir, IndexFile)
		if fileExists(dirIndexPath) {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("UP"))
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DOWN"))
	}
}
