package main

import (
	"strings"
	"testing"

	"github.com/mauhlik/go-index/config"
	"github.com/mauhlik/go-index/internal/go-index/providers"
	"github.com/sirupsen/logrus"
)

func TestSetupProviderForRepository(t *testing.T) {
	t.Parallel()

	logger := logrus.New()
	tests := getTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := setupProviderForRepository(testCase.cfg, testCase.repo, logger)
			handleError(t, err, testCase.wantErrPart)

			if err != nil {
				return
			}

			verifyProviderType(t, got, testCase.wantType)
		})
	}
}

func handleError(t *testing.T, err error, wantErrPart string) {
	t.Helper()

	if wantErrPart != "" {
		if err == nil {
			t.Fatalf("expected error containing %q", wantErrPart)
		}

		if !strings.Contains(err.Error(), wantErrPart) {
			t.Fatalf("error %q does not contain %q", err.Error(), wantErrPart)
		}
	}
}

func verifyProviderType(t *testing.T, got interface{}, wantType string) {
	t.Helper()

	switch wantType {
	case "local":
		_, ok := got.(*providers.LocalProvider)
		if !ok {
			t.Errorf("expected *LocalProvider, got %T", got)
		}
	case "s3":
		_, ok := got.(*providers.S3Provider)
		if !ok {
			t.Errorf("expected *S3Provider, got %T", got)
		}
	}
}

//nolint:funlen
func getTestCases() []struct {
	name        string
	cfg         *config.Config
	repo        config.RepositoryConfig
	wantType    string
	wantErrPart string
} {
	return []struct {
		name        string
		cfg         *config.Config
		repo        config.RepositoryConfig
		wantType    string
		wantErrPart string
	}{
		{
			name: "provider not found",
			cfg: &config.Config{
				Port:         "8080",
				Repositories: []config.RepositoryConfig{},
				Providers:    map[string]interface{}{},
			},
			repo: config.RepositoryConfig{
				Name:     "repo1",
				Provider: "nonexistent",
			},
			wantType:    "",
			wantErrPart: "provider not found",
		},
		{
			name: "unknown provider type",
			cfg: &config.Config{
				Port:         "8080",
				Repositories: []config.RepositoryConfig{},
				Providers: map[string]interface{}{
					"unknown": "invalid-type",
				},
			},
			repo: config.RepositoryConfig{
				Name:     "repo2",
				Provider: "unknown",
			},
			wantType:    "",
			wantErrPart: "unknown provider type",
		},
		{
			name: "local provider success",
			cfg: &config.Config{
				Port:         "8080",
				Repositories: []config.RepositoryConfig{},
				Providers: map[string]interface{}{
					"local": config.LocalProviderConfig{
						Path: "some/path",
						Type: "local",
					},
				},
			},
			repo: config.RepositoryConfig{
				Name:     "repo3",
				Provider: "local",
			},
			wantType:    "local",
			wantErrPart: "",
		},
		{
			name: "s3 provider success",
			cfg: &config.Config{
				Port:         "8080",
				Repositories: []config.RepositoryConfig{},
				Providers: map[string]interface{}{
					"s3": config.S3ProviderConfig{
						Bucket:    "bucket",
						Endpoint:  "https://s3.example.com",
						AccessKey: "access",
						SecretKey: "secret",
						Region:    "region",
						Type:      "s3",
					},
				},
			},
			repo: config.RepositoryConfig{
				Name:     "repo4",
				Provider: "s3",
			},
			wantType:    "s3",
			wantErrPart: "",
		},
	}
}
