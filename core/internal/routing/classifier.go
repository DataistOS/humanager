// core/internal/routing/classifier.go
package routing

import (
	"path/filepath"
	"strings"
)

type Classifier struct {
}

func NewClassifier() *Classifier {
	return &Classifier{}
}

func (c *Classifier) RouteFile(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".jpg", ".png", ".jpeg", ".gif", ".webp":
		return filepath.Join("nervous_system", "sensory_system", "visual")

	case ".mp3", ".wav", ".aac", ".flac":
		return filepath.Join("nervous_system", "sensory_system", "Auditory system")

	case ".md", ".txt", ".pdf", ".docx":
		return filepath.Join("nervous_system", "central_nervous_system", "brain", "mind", "memory", "SQ3R")

	case ".mp4", ".mkv", ".mov":
		return filepath.Join("nervous_system", "sensory_system", "visual")

	default:
		return filepath.Join("digestive_system")
	}
}
