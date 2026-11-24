// core/internal/api/handler.go
package api

import (
	"encoding/json"
	"humanager/internal/routing"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DirItem struct {
	Name   string `json:"name"`
	IsDirs bool   `json:"is_dir"`
	Path   string `json:"path"`
}

func resolveDirPath(storageRoot, relPath string) (string, string) {
	cleanedRel := filepath.ToSlash(filepath.Clean("/" + strings.TrimSpace(relPath)))
	cleanedRel = strings.TrimPrefix(cleanedRel, "/")
	if cleanedRel == "." {
		cleanedRel = ""
	}

	candidate := filepath.Join(storageRoot, cleanedRel)
	if _, err := os.Stat(candidate); err == nil {
		return candidate, cleanedRel
	}

	if cleanedRel != "" && !strings.HasPrefix(cleanedRel, "biological_system") {
		bioRel := filepath.ToSlash(filepath.Join("biological_system", cleanedRel))
		bioCandidate := filepath.Join(storageRoot, bioRel)
		if _, err := os.Stat(bioCandidate); err == nil {
			return bioCandidate, bioRel
		}
	}

	return candidate, cleanedRel
}

func HandleGetTree(storageRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		systems := make([]string, 0)
		bioDir := filepath.Join(storageRoot, "biological_system")
		entries, err := os.ReadDir(bioDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					systems = append(systems, entry.Name())
				}
			}
		}

		json.NewEncoder(w).Encode(map[string]any{
			"root":    filepath.Base(filepath.Clean(storageRoot)),
			"systems": systems,
		})
	}
}

func HandleListDirectory(storageRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		relPath := r.URL.Query().Get("path")
		targetPath, effectiveRel := resolveDirPath(storageRoot, relPath)

		entries, err := os.ReadDir(targetPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]DirItem, 0, len(entries))
		for _, entry := range entries {
			itemRelPath := filepath.ToSlash(filepath.Join(effectiveRel, entry.Name()))
			items = append(items, DirItem{
				Name:   entry.Name(),
				IsDirs: entry.IsDir(),
				Path:   itemRelPath,
			})
		}

		json.NewEncoder(w).Encode(map[string]any{
			"current_path": effectiveRel,
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

		collection := r.FormValue("collection")
		slug := r.FormValue("slug")
		if slug == "" {
			slug = r.FormValue("finglish")
		}
		recordID := r.FormValue("record_id")

		var relativeSubPath string
		var finalFileName string

		if routedDir, routedFileName, matched := classifier.RouteCollectionRecord(collection, slug, recordID, header.Filename); matched {
			relativeSubPath = routedDir
			finalFileName = routedFileName
		} else {
			relativeSubPath = classifier.RouteFile(header.Filename)
			finalFileName = header.Filename
		}

		targetDir := filepath.Join(storageRoot, relativeSubPath)

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			http.Error(w, "Failed to create target directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		dstPath := filepath.Join(targetDir, finalFileName)
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

		relativePath := filepath.ToSlash(filepath.Join(relativeSubPath, finalFileName))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success":       true,
			"message":       "File uploaded and classified successfully",
			"assigned_path": filepath.ToSlash(relativeSubPath),
			"relative_path": relativePath,
			"file_name":     finalFileName,
		})
	}
}
