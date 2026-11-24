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

	rootPath := filepath.Clean(basePath)
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return fmt.Errorf("failed to create root path: %v", err)
	}

	for _, branch := range anatomy.Branches {
		systemPath := filepath.Join(rootPath, branch.System)
		if err := os.MkdirAll(systemPath, 0755); err != nil {
			return fmt.Errorf("failed to create system path %s: %v", systemPath, err)
		}

		for _, organ := range branch.Organs {
			organPath := filepath.Join(systemPath, organ)
			if err := os.MkdirAll(organPath, 0755); err != nil {
				return fmt.Errorf("failed to create organ path %s: %v", organPath, err)
			}
		}

		for _, sub := range branch.Subsystems {
			if err := buildSubsystem(systemPath, sub); err != nil {
				return err
			}
		}
	}

	fmt.Println("✨ Anatomy tree successfully built on disk at:", rootPath)
	return nil
}

func buildSubsystem(parentPath string, sub Subsystem) error {
	currentPath := filepath.Join(parentPath, sub.Name)
	if err := os.MkdirAll(currentPath, 0755); err != nil {
		return fmt.Errorf("failed to create subsystem path %s: %v", currentPath, err)
	}

	for _, organ := range sub.Organs {
		if err := os.MkdirAll(filepath.Join(currentPath, organ), 0755); err != nil {
			return fmt.Errorf("failed to create organ path %s: %v", organ, err)
		}
	}

	for _, cat := range sub.Categories {
		if err := os.MkdirAll(filepath.Join(currentPath, cat), 0755); err != nil {
			return fmt.Errorf("failed to create category path %s: %v", cat, err)
		}
	}

	for _, childSub := range sub.Subsystems {
		if err := buildSubsystem(currentPath, childSub); err != nil {
			return err
		}
	}

	return nil
}
