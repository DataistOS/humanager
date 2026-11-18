// core/internal/engine/tree_builder.go
package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AnatomyNode struct {
	RootName string   `json:"root_name"`
	Branches []Branch `json:"branches"`
}

type Branch struct {
	System     string      `json:"system"`
	Organs     []string    `json:"organs,omitempty"`
	Subsystems []Subsystem `json:"subsystems,omitempty"`
}

type Subsystem struct {
	Name       string      `json:"name"`
	Organs     []string    `json:"organs,omitempty"`
	Categories []string    `json:"categories,omitempty"`
	Subsystems []Subsystem `json:"subsystems,omitempty"`
}

func BuildAnatomyTree(basePath string, jsonConfigPath string) error {
	fileData, err := os.ReadFile(jsonConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read anatomy json config: %v", err)
	}

	var anatomy AnatomyNode
	if err := json.Unmarshal(fileData, &anatomy); err != nil {
		return fmt.Errorf("failed to parse anatomy json: %v", err)
	}

	rootPath := filepath.Join(basePath, anatomy.RootName)
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return fmt.Errorf("failed to create root path: %v", err)
	}

	for _, branch := range anatomy.Branches {
		systemPath := filepath.Join(rootPath, branch.System)

		for _, organ := range branch.Organs {
			organPath := filepath.Join(systemPath, organ)
			os.MkdirAll(organPath, 0755)
		}

		for _, sub := range branch.Subsystems {
			buildSubsystem(systemPath, sub)
		}
	}

	fmt.Println("✨ Anatomy tree successfully built on disk!")
	return nil
}

func buildSubsystem(parentPath string, sub Subsystem) {
	currentPath := filepath.Join(parentPath, sub.Name)
	os.MkdirAll(currentPath, 0755)

	for _, organ := range sub.Organs {
		os.MkdirAll(filepath.Join(currentPath, organ), 0755)
	}

	for _, cat := range sub.Categories {
		os.MkdirAll(filepath.Join(currentPath, cat), 0755)
	}

	for _, childSub := range sub.Subsystems {
		buildSubsystem(currentPath, childSub)
	}
}
