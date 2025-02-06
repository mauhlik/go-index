package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
)

type S3Client interface {
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input,
		opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type S3Provider struct {
	Client S3Client
	Bucket string
	logger *logrus.Logger
}

func NewS3Provider(bucket, endpoint, accessKey, secretKey, region string, logger *logrus.Logger) (*S3Provider, error) {
	logger.Infof("Initialized S3 client endpoint %s region %s", endpoint, region)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithBaseEndpoint(endpoint),
	)

	if err != nil {
		logger.WithError(err).Error("Failed to load AWS config")

		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	logger.Infof("Initialized S3 client endpoint %s region %s", endpoint, region)

	return &S3Provider{Client: client, Bucket: bucket, logger: logger}, nil
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
			p.logger.WithError(err).Error("Failed to list objects")

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
