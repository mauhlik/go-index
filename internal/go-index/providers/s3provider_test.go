package providers_test

import (
	"testing"

	"github.com/MaUhlik-cen56998/go-index/internal/go-index/mocks"
	"github.com/MaUhlik-cen56998/go-index/internal/go-index/providers"
	"github.com/golang/mock/gomock"
)

func TestS3ProviderGetVersions(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3Client := mocks.NewMockS3Client(mockCtrl)
	provider := &providers.S3Provider{
		Client: mockS3Client,
		Bucket: "test-bucket",
		Logger: nil,
	}

	moduleName := "fe"
	artifactName := "app1"
	expectedVersions := []string{"0.0.0", "0.0.1", "1.0.0", "2.0.0"}

	gotVersions, err := provider.GetVersions(moduleName, artifactName)
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
