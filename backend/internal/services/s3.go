package services

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	bucketName string
	client     *s3.Client
}

func NewS3Service(bucketName string, client *s3.Client) *S3Service {
	return &S3Service{
		bucketName: bucketName,
		client:     client,
	}
}

func (s *S3Service) PutObject(ctx context.Context, key string, data io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(s.bucketName),
		Key:          aws.String(key),
		Body:         data,
		CacheControl: aws.String("public, max-age=2592000"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Service) GetObject(ctx context.Context, key string) error {
	// Not necessary due to using concatenation instead
	// `${AWS_CDN_URL}/${AWS_BUCKET_NAME}/${key}`
	return nil
}
