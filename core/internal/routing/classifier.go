// core/internal/routing/classifier.go
package routing

import (
	"path/filepath"
	"strings"
	"unicode"
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

func (c *Classifier) RouteCollectionRecord(collection, slugOrFinglish, recordID, originalFileName string) (dirRelPath string, fileName string, ok bool) {
	if strings.EqualFold(strings.TrimSpace(collection), "daftarche") {
		cleanSlug := sanitizeSlug(slugOrFinglish)
		if cleanSlug == "" {
			cleanSlug = sanitizeSlug(strings.TrimSuffix(originalFileName, filepath.Ext(originalFileName)))
		}
		if cleanSlug == "" {
			cleanSlug = "contact"
		}

		alpha := extractAlphaFolder(cleanSlug)
		folderName := cleanSlug
		if strings.TrimSpace(recordID) != "" {
			folderName = cleanSlug + "__" + strings.TrimSpace(recordID)
		}

		baseDaftarche := filepath.Join(
			"biological_system",
			"nervous_system",
			"central_nervous_system",
			"brain",
			"mind",
			"memory",
			"socialization",
			"daftarche",
		)

		return filepath.Join(baseDaftarche, alpha, folderName), "avatar.jpg", true
	}

	return "", "", false
}

func extractAlphaFolder(s string) string {
	for _, r := range strings.TrimSpace(s) {
		if unicode.IsLetter(r) && r <= unicode.MaxASCII {
			return strings.ToUpper(string(r))
		}
	}
	return "#"
}

func sanitizeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")

	var b strings.Builder
	prevUnderscore := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevUnderscore = false
		} else if r == '_' && !prevUnderscore && b.Len() > 0 {
			b.WriteRune('_')
			prevUnderscore = true
		}
	}

	return strings.Trim(b.String(), "_")
}
