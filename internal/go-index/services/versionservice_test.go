package services_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mauhlik/go-index/internal/go-index/mocks"
	"github.com/mauhlik/go-index/internal/go-index/services"
	"github.com/sirupsen/logrus"
)

func TestVersionServiceGetVersions(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProvider := mocks.NewMockProvider(mockCtrl)
	logger := logrus.New()
	service := services.NewService(mockProvider, logger)

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
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProvider := mocks.NewMockProvider(mockCtrl)
	logger := logrus.New()
	service := services.NewService(mockProvider, logger)

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
