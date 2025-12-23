package test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"luny.dev/cherryauctions/services"
)

func TestS3Buckets(t *testing.T) {
	s3Client := services.NewS3Service()

	ctx := t.Context()
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String("test-bucket"),
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String("test-bucket")})
	if err != nil {
		t.Fatal(err.Error())
	}
}
