package providers

import (
	"fmt"
	"os"
	"path/filepath"
)

type LocalProvider struct {
	basePath string
}

func NewLocalProvider(basePath string) *LocalProvider {
	return &LocalProvider{basePath: basePath}
}

func (p *LocalProvider) GetVersions(moduleName, artifactName string) ([]string, error) {
	path := filepath.Join(p.basePath, moduleName, artifactName)
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var versions []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		version := ExtractVersionFromFilename(filename, artifactName)

		if version != "" {
			versions = append(versions, version)
		}
	}

	return versions, nil
}
