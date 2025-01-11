package services

import (
	"testing"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
	"github.com/golang/mock/gomock"
)

type MockProvider struct {
	providers.Provider
	mock *gomock.Controller
}

func (m *MockProvider) GetVersions(moduleName, artifactName string) ([]string, error) {
	return []string{"0.0.0", "0.0.1", "1.0.0", "2.0.0"}, nil
}

func TestVersionServiceGetVersions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProvider := &MockProvider{mock: mockCtrl}
	service := NewService(mockProvider)

	moduleName := "fe"
	artifactName := "app1"
	expectedVersions := []string{"0.0.0", "0.0.1", "1.0.0", "2.0.0"}

	gotVersions, err := service.GetVersions(moduleName, artifactName)
	if err != nil {
		t.Fatalf("GetVersions returned an error: %v", err)
	}

	if len(gotVersions) != len(expectedVersions) {
		t.Errorf("GetVersions returned %d versions; want %d", len(gotVersions), len(expectedVersions))
	}
	for i, version := range gotVersions {
		if version != expectedVersions[i] {
			t.Errorf("GetVersions returned version %q; want %q", version, expectedVersions[i])
		}
	}
}

func TestVersionServiceGetLatestVersion(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProvider := &MockProvider{mock: mockCtrl}
	service := NewService(mockProvider)

	moduleName := "fe"
	artifactName := "app1"
	expectedLatestVersion := "2.0.0"

	gotLatestVersion, err := service.GetLatestVersion(moduleName, artifactName)
	if err != nil {
		t.Fatalf("GetLatestVersion returned an error: %v", err)
	}

	if gotLatestVersion != expectedLatestVersion {
		t.Errorf("GetLatestVersion returned %q; want %q", gotLatestVersion, expectedLatestVersion)
	}
}
