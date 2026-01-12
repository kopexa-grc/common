// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package s3store

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kopexa-grc/common/blob/driver"
	goblob "gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

const S3_DEFAULT_REGION = "de"

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string
	BucketPrefix    string
	UsePathStyle    bool
	ContainerName   string
}

type S3Service struct {
	client        *s3.Client
	bucketPrefix  string
	containerName string
	config        *S3Config
}

func NewS3Service(cfg *S3Config) (*S3Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("s3store: config cannot be nil")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("s3store: failed to load AWS config: %w", err)
	}

	clientOptions := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = cfg.UsePathStyle
		},
	}

	if cfg.Endpoint != "" {
		clientOptions = append(clientOptions, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOptions...)

	service := &S3Service{
		client:        client,
		containerName: cfg.ContainerName,
		bucketPrefix:  cfg.BucketPrefix,
		config:        cfg,
	}

	return service, nil
}

func (s *S3Service) Upload(ctx context.Context, key string, reader io.Reader, metadata map[string]string) error {
	if key == "" {
		return fmt.Errorf("s3store: key cannot be empty")
	}

	bucket, err := s3blob.OpenBucketV2(ctx, s.client, s.bucketPrefix, nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	w, err := bucket.NewWriter(ctx, key, nil)
	if err != nil {
		return err
	}

	_, writeErr := io.Copy(w, reader)
	// Always check the return value of Close when writing.
	closeErr := w.Close()
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	if closeErr != nil {
		log.Fatal(closeErr)
	}

	return nil
}

// UploadTyped lädt ein Objekt mit Content-Type und optionalen Headern/Metadaten.
func (s *S3Service) UploadTyped(ctx context.Context, key string, reader io.Reader, contentType string, opts *driver.WriterOptions) error {
	if key == "" {
		return fmt.Errorf("s3store: key cannot be empty")
	}

	bucket, err := s3blob.OpenBucketV2(ctx, s.client, s.bucketPrefix, nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	w, err := bucket.NewWriter(ctx, key, &goblob.WriterOptions{
		ContentType: contentType,
		Metadata:    opts.Metadata,
	})
	if err != nil {
		return err
	}

	_, writeErr := io.Copy(w, reader)
	// Always check the return value of Close when writing.
	closeErr := w.Close()
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	if closeErr != nil {
		log.Fatal(closeErr)
	}

	return nil
}

func (s *S3Service) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	if key == "" {
		return nil, fmt.Errorf("s3store: key cannot be empty")
	}

	bucket, err := s3blob.OpenBucketV2(ctx, s.client, s.bucketPrefix, nil)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	r, err := bucket.NewReader(ctx, key, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *S3Service) RangeDownload(ctx context.Context, key string, offset, length int64) (*goblob.Reader, error) {
	if key == "" {
		return nil, fmt.Errorf("s3store: key cannot be empty")
	}

	bucket, err := s3blob.OpenBucketV2(ctx, s.client, s.bucketPrefix, nil)
	if err != nil {
		return nil, err
	}
	defer bucket.Close()

	reader, err := bucket.NewRangeReader(ctx, key, offset, length, nil)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (s *S3Service) Delete(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("s3store: key cannot be empty")
	}

	bucket, err := s3blob.OpenBucketV2(ctx, s.client, s.bucketPrefix, nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	return bucket.Delete(ctx, key)
}

type CopyParams struct {
	SourceKey string
	DestKey   string
}

func (s *S3Service) Copy(ctx context.Context, params CopyParams) error {
	reader, err := s.Download(ctx, params.SourceKey)
	if err != nil {
		return fmt.Errorf("s3store: failed to download source object: %w", err)
	}
	defer reader.Close()

	err = s.Upload(ctx, params.DestKey, reader, nil)
	if err != nil {
		return fmt.Errorf("s3store: failed to upload copied object: %w", err)
	}

	return nil

	// does not work
	// trimmedBucketName := strings.TrimSuffix(s.bucketPrefix, "/")
	// trimmedBucketName := strings.TrimSuffix(s.bucketPrefix, "/")

	// sourceKey := fmt.Sprintf("%s/%s/%s", trimmedBucketName, trimmedBucketName, params.SourceKey)
	// destKey := fmt.Sprintf("%s/%s", trimmedBucketName, params.DestKey)

	// fmt.Printf("source: %s \n", sourceKey)
	// fmt.Printf("dest: %s \n", destKey)

	// _, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
	// 	Bucket: aws.String(trimmedBucketName),
	// 	Key:    aws.String(destKey),

	// 	CopySource: aws.String(sourceKey),
	// })
	// if err != nil {
	// 	return fmt.Errorf("s3store: failed to copy object: %w", err)
	// }
}

func (s *S3Service) GetSignedURL(ctx context.Context, key string, expiration time.Duration, method string) (string, error) {
	presignClient := s3.NewPresignClient(s.client, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	var output *v4.PresignedHTTPRequest
	var err error

	switch method {
	case "PUT":
		output, err = presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucketPrefix),
			Key:    aws.String(key),
		})
	case "DELETE":
		output, err = presignClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketPrefix),
			Key:    aws.String(key),
		})
	default:
		output, err = presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucketPrefix),
			Key:    aws.String(key),
		})
	}

	if err != nil {
		return "", fmt.Errorf("s3store: failed to generate signed URL: %w", err)
	}

	return output.URL, nil
}
