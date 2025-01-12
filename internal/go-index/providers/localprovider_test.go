package providers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
)

func TestLocalProviderGetVersions(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	defer os.RemoveAll(tempDir)

	moduleName := "fe"
	artifactName := "app1"
	versions := []string{"0.0.0", "0.0.1", "1.0.0", "2.0.0"}

	for _, version := range versions {
		filename := filepath.Join(tempDir, moduleName, artifactName, artifactName+"-"+version+".txt")
		err := os.MkdirAll(filepath.Dir(filename), 0755)

		if err != nil {
			t.Fatalf("Failed to create directories: %v", err)
		}

		_, err = os.Create(filename)

		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	provider := providers.NewLocalProvider(tempDir)

	gotVersions, err := provider.GetVersions(moduleName, artifactName)
	if err != nil {
		t.Fatalf("GetVersions returned an error: %v", err)
	}

	expectedVersions := []string{"0.0.0", "0.0.1", "1.0.0", "2.0.0"}
	if len(gotVersions) != len(expectedVersions) {
		t.Errorf("GetVersions returned %d versions; want %d", len(gotVersions), len(expectedVersions))
	}

	for i, version := range gotVersions {
		if version != expectedVersions[i] {
			t.Errorf("GetVersions returned version %q; want %q", version, expectedVersions[i])
		}
	}
}
