package lfe

import (
	"fmt"
	"os"
	"path/filepath"
)

type discoveredLicense struct{ path, json string }

func discoverLicense(explicit string) (discoveredLicense, error) {
	if explicit != "" {
		return readLicense(explicit)
	}
	paths := []string{"license.json"}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".lfe-be", "license.json"))
	}
	paths = append(paths, "/etc/lfe-be/license.json")
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return readLicense(p)
		}
	}
	return discoveredLicense{}, fmt.Errorf("lfe: license required; searched %v", paths)
}

func readLicense(path string) (discoveredLicense, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return discoveredLicense{}, fmt.Errorf("lfe: read license %q: %w", path, err)
	}
	return discoveredLicense{path: path, json: string(b)}, nil
}
