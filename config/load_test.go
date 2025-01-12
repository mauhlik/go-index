package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MaUhlik-cen56998/go-index/config"
)

func TestLoadConfigJSON(t *testing.T) {
	t.Parallel()

	configContent := `
{
    "repositories": [
        {
            "name": "localrepo1",
            "provider": "test-local-provider"
        },
        {
            "name": "s3repo1",
            "provider": "s3"
        }
    ],
    "providers": {
        "test-local-provider": {
            "type": "local",
            "path": "test-local-provider"
        },
        "s3": {
            "type": "s3",
            "bucket": "test-bucket",
            "endpoint": "https://s3.com",
            "accessKey": "9wnheHR37PwXdE8YF56U",
            "secretKey": "MzTDkSsvcJgHAoSU2D4Z7",
            "region": "main"
        }
    }
}
`
	testLoadConfig(t, configContent, ".json")
}

func TestLoadConfigYAML(t *testing.T) {
	t.Parallel()

	configContent := `
repositories:
  - name: localrepo1
    provider: test-local-provider
  - name: s3repo1
    provider: s3
providers:
  test-local-provider:
    type: local
    path: test-local-provider
  s3:
    type: s3
    bucket: test-bucket
    endpoint: https://s3.com
    accessKey: 9wnheHR37PwXdE8YF56U
    secretKey: MzTDkSsvcJgHAoSU2D4Z7
    region: main
`
	testLoadConfig(t, configContent, ".yaml")
}

func testLoadConfig(t *testing.T, configContent, ext string) {
	t.Helper()
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config"+ext)

	err := os.WriteFile(configFile, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	checkRepositories(t, cfg.Repositories)
	checkProviders(t, cfg.Providers)
}

func checkRepositories(t *testing.T, repositories []config.RepositoryConfig) {
	t.Helper()

	if len(repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(repositories))
	}
}

func checkProviders(t *testing.T, providers map[string]interface{}) {
	t.Helper()

	if len(providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(providers))
	}

	checkLocalProvider(t, providers["test-local-provider"])
	checkS3Provider(t, providers["s3"])
}

func checkLocalProvider(t *testing.T, provider interface{}) {
	t.Helper()

	localProvider, ok := provider.(config.LocalProviderConfig)
	if !ok {
		t.Fatalf("Expected LocalProviderConfig, got %T", provider)
	}

	if localProvider.Path != "test-local-provider" {
		t.Errorf("Expected path 'test-local-provider', got '%s'", localProvider.Path)
	}
}

func checkS3Provider(t *testing.T, provider interface{}) {
	t.Helper()

	s3Provider, ok := provider.(config.S3ProviderConfig)
	if !ok {
		t.Fatalf("Expected S3ProviderConfig, got %T", provider)
	}

	if s3Provider.Bucket != "test-bucket" {
		t.Errorf("Expected bucket 'test-bucket', got '%s'", s3Provider.Bucket)
	}

	if s3Provider.Endpoint != "https://s3.com" {
		t.Errorf("Expected endpoint 'https://s3.com', got '%s'", s3Provider.Endpoint)
	}

	if s3Provider.AccessKey != "9wnheHR37PwXdE8YF56U" {
		t.Errorf("Expected accessKey '9wnheHR37PwXdE8YF56U', got '%s'", s3Provider.AccessKey)
	}

	if s3Provider.SecretKey != "MzTDkSsvcJgHAoSU2D4Z7" {
		t.Errorf("Expected secretKey 'MzTDkSsvcJgHAoSU2D4Z7', got '%s'", s3Provider.SecretKey)
	}

	if s3Provider.Region != "main" {
		t.Errorf("Expected region 'main', got '%s'", s3Provider.Region)
	}
}
