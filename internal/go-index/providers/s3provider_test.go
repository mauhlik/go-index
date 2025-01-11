package providers

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
)

type MockS3Client struct {
	s3.Client
	mock *gomock.Controller
}

func (m *MockS3Client) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	output := &s3.ListObjectsV2Output{
		Contents: []types.Object{
			{Key: aws.String("fe/app1/app1-0.0.0.txt")},
			{Key: aws.String("fe/app1/app1-0.0.1.txt")},
			{Key: aws.String("fe/app1/app1-1.0.0.txt")},
			{Key: aws.String("fe/app1/app1-2.0.0.txt")},
		},
	}
	return output, nil
}

func TestS3ProviderGetVersions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockS3Client := &MockS3Client{mock: mockCtrl}
	provider := &S3Provider{
		client: mockS3Client,
		bucket: "test-bucket",
		logger: nil,
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
