package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"gofr.dev/pkg/gofr/logging"
)

type S3Client interface {
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input,
		opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type S3Provider struct {
	Client S3Client
	Bucket string
	Logger logging.Logger
}

func NewS3Provider(bucket, endpoint, accessKey, secretKey, region string, logger logging.Logger) (*S3Provider, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithBaseEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	logger.Infof("Initialized S3 provider with bucket: %s, endpoint: %s", bucket, endpoint)

	return &S3Provider{Client: client, Bucket: bucket, Logger: logger}, nil
}

func (p *S3Provider) GetVersions(moduleName, artifactName string) ([]string, error) {
	prefix := fmt.Sprintf("%s/%s/", moduleName, artifactName)
	input := &s3.ListObjectsV2Input{
		Bucket:                   &p.Bucket,
		Prefix:                   &prefix,
		ContinuationToken:        nil,
		Delimiter:                nil,
		EncodingType:             "",
		ExpectedBucketOwner:      nil,
		FetchOwner:               aws.Bool(false),
		MaxKeys:                  aws.Int32(0),
		OptionalObjectAttributes: nil,
		RequestPayer:             "",
		StartAfter:               aws.String(""),
	}

	var versions []string

	paginator := s3.NewListObjectsV2Paginator(p.Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			key := *obj.Key
			if strings.HasPrefix(key, prefix) {
				filename := strings.TrimPrefix(key, prefix)
				version := ExtractVersionFromFilename(filename, artifactName)

				if version != "" {
					versions = append(versions, version)
				}
			}
		}
	}

	return versions, nil
}
