// core/internal/api/router.go
package api

import (
	"encoding/json"
	"net/http"
)

func SetupRouter(storageRoot string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", handleHealthCheck)

	mux.HandleFunc("/api/tree", handleGetAnatomyTree)

	mux.HandleFunc("/api/dir", HandleListDirectory(storageRoot))

	mux.HandleFunc("/api/mkdir", HandleCreateDirectory(storageRoot))

	mux.HandleFunc("/api/upload", HandleUploadAndClassify(storageRoot))

	return mux
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "running",
		"message": "HuMaNager Core Daemon is active!",
	})
}

func handleGetAnatomyTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"root": "secondbrain_storage",
		"systems": []string{
			"digestive_system",
			"genital_system",
			"nervous_system",
		},
	})
}
