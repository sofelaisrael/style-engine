package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func stylesDir() string {
	// Try to find styles/ relative to the project root
	// First try: relative to executable
	if ex, err := os.Executable(); err == nil {
		dir := filepath.Dir(ex)
		candidate := filepath.Join(dir, "styles")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		// Try one level up (for go run temp dir case)
		candidate = filepath.Join(dir, "..", "styles")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	// Second try: look for a marker file (go.mod) walking up from CWD
	dir, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			candidate := filepath.Join(dir, "styles")
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				return candidate
			}
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback: just use CWD
	return "styles"
}

func LoadProfile(name string) (*Profile, error) {
	data, err := os.ReadFile(filepath.Join(stylesDir(), name+".json"))
	if err != nil {
		return nil, fmt.Errorf("style %q not found: %w", name, err)
	}

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse style %q: %w", name, err)
	}
	return &profile, nil
}

func ListStyles() ([]string, error) {
	entries, err := os.ReadDir(stylesDir())
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
