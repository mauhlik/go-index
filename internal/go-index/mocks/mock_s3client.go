package mocks

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go/middleware"
	"github.com/golang/mock/gomock"
)

const restoreExpiryHours = 24
const size = 123

type MockS3Client struct {
	s3.Client
	mock *gomock.Controller
}

func NewMockS3Client(ctrl *gomock.Controller) *MockS3Client {
	return &MockS3Client{
		Client: s3.Client{},
		mock:   ctrl,
	}
}

func (m *MockS3Client) ListObjectsV2(_ context.Context, _ *s3.ListObjectsV2Input,
	_ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	objects := []types.Object{}
	keys := []string{"fe/app1/app1-0.0.0.txt", "fe/app1/app1-0.0.1.txt",
		"fe/app1/app1-1.0.0.txt", "fe/app1/app1-2.0.0.txt"}
	etags := []string{"etag-0.0.0", "etag-0.0.1", "etag-1.0.0", "etag-2.0.0"}

	for i, key := range keys {
		objects = append(objects, types.Object{
			Key:               aws.String(key),
			LastModified:      aws.Time(time.Now()),
			ETag:              aws.String(etags[i]),
			Size:              aws.Int64(size),
			StorageClass:      types.ObjectStorageClassStandard,
			ChecksumAlgorithm: []types.ChecksumAlgorithm{},
			Owner:             &types.Owner{ID: aws.String("owner-id"), DisplayName: aws.String("owner-name")},
			RestoreStatus: &types.RestoreStatus{
				IsRestoreInProgress: aws.Bool(false),
				RestoreExpiryDate:   aws.Time(time.Now().Add(restoreExpiryHours * time.Hour)),
			},
		})
	}

	output := &s3.ListObjectsV2Output{
		CommonPrefixes:        []types.CommonPrefix{},
		ContinuationToken:     aws.String(""),
		Delimiter:             aws.String(""),
		EncodingType:          types.EncodingTypeUrl,
		IsTruncated:           aws.Bool(false),
		KeyCount:              aws.Int32(0),
		MaxKeys:               aws.Int32(0),
		Name:                  aws.String(""),
		NextContinuationToken: aws.String(""),
		Prefix:                aws.String(""),
		RequestCharged:        types.RequestChargedRequester,
		StartAfter:            aws.String(""),
		ResultMetadata:        middleware.Metadata{},
		Contents:              objects,
	}

	return output, nil
}
