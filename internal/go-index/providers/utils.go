package providers

import (
	"path/filepath"
	"strings"
)

func ExtractVersionFromFilename(filename, artifactName string) string {
	if strings.HasPrefix(filename, artifactName+"-") {
		version := strings.TrimPrefix(filename, artifactName+"-")

		for {
			ext := filepath.Ext(version)
			if ext == "" {
				break
			}

			if ContainsNumbers(ext) {
				break
			}
			version = strings.TrimSuffix(version, ext)
		}

		return version
	}

	return ""
}

func ContainsNumbers(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}

	return false
}
