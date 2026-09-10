package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadProfile(name string) (*Profile, error) {
	data, err := os.ReadFile(filepath.Join("styles", name+".json"))
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
	entries, err := os.ReadDir("styles")
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
