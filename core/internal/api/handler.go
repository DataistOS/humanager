// core/internal/api/handler.go
package api

import (
	"encoding/json"
	"humanager/internal/routing"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type DirItem struct {
	Name   string `json:"name"`
	IsDirs bool   `json:"is_dir"`
	Path   string `json:"path"`
}

func HandleListDirectory(storageRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		relPath := r.URL.Query().Get("path")
		targetPath := filepath.Join(storageRoot, relPath)

		entries, err := os.ReadDir(targetPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var items []DirItem
		for _, entry := range entries {
			items = append(items, DirItem{
				Name:   entry.Name(),
				IsDirs: entry.IsDir(),
				Path:   filepath.Join(relPath, entry.Name()),
			})
		}

		json.NewEncoder(w).Encode(map[string]any{
			"current_path": relPath,
			"items":        items,
		})
	}
}

func HandleCreateDirectory(storageRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Path string `json:"path"`
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		targetPath := filepath.Join(storageRoot, req.Path, req.Name)

		err := os.MkdirAll(targetPath, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "Folder created successfully",
			"path":    targetPath,
		})
	}
}

func HandleUploadAndClassify(storageRoot string) http.HandlerFunc {
	classifier := routing.NewClassifier()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		relativeSubPath := classifier.RouteFile(header.Filename)
		targetDir := filepath.Join(storageRoot, relativeSubPath)

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			http.Error(w, "Failed to create target directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		dstPath := filepath.Join(targetDir, header.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Failed to create destination file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success":       true,
			"message":       "File uploaded and classified successfully",
			"assigned_path": relativeSubPath,
			"file_name":     header.Filename,
		})
	}
}
