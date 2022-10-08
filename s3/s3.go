package s3

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Manager struct {
	ctx        context.Context
	downloader *manager.Downloader
}

func AWSConfigForTest(ctx context.Context) aws.Config {
	const region = "us-east-1"
	endpoint := "http://localhost:9000"
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})

	creds := credentials.NewStaticCredentialsProvider(
		"minio", // these values match the credentials of bitnami/minio docker image
		"miniosecret",
		"",
	)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(resolver),
		config.WithClientLogMode(aws.LogRetries),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewS3Manager(ctx context.Context) *S3Manager {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithClientLogMode(aws.LogRetries))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	// awsConfig := AWSConfigForTest(ctx)
	client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true // to work with minio
	})

	manager := &S3Manager{
		ctx:        ctx,
		downloader: manager.NewDownloader(client),
	}
	manager.downloader.Concurrency = 1

	return manager
}

func (s3m *S3Manager) GetS3Reader(bucket string, key string, maxBufSize int) io.Reader {
	buf := NewWriteAtReader(maxBufSize)
	go func() {
		input := &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		}
		_, err := s3m.downloader.Download(s3m.ctx, buf, input)
		if err != nil {
			log.Fatalf("failed to download s3 file: %v", err)
		}
		buf.Done()
	}()

	return buf
}
