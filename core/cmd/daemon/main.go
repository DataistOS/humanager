// core/cmd/daemon/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"humanager/internal/api"
	"humanager/internal/engine"
)

func resolveAnatomyConfigPath() string {
	candidates := []string{
		"configs/anatomy_tree.json",
		"core/configs/anatomy_tree.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return candidates[0]
}

func main() {
	fmt.Println("🧠 HuMaNager Daemon is starting...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	storageRoot := filepath.Join(homeDir, "5ec0ndbr41n")
	if err := os.MkdirAll(storageRoot, 0755); err != nil {
		log.Fatalf("Failed to create storage root: %v", err)
	}

	jsonConfigPath := resolveAnatomyConfigPath()
	if err := engine.BuildAnatomyTree(storageRoot, jsonConfigPath); err != nil {
		log.Printf("Warning: Could not build anatomy tree from JSON (%s): %v", jsonConfigPath, err)
	}

	router := api.SetupRouter(storageRoot)
	port := ":8081"

	fmt.Printf("🚀 API Server is running on http://localhost%s (root: %s)\n", port, storageRoot)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
