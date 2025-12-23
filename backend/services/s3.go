package services

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"luny.dev/cherryauctions/utils"
)

func NewS3Service() *s3.Client {
	// Load the Shared AWS Configuration.
	// Need some environment variables:
	// AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN.
	//
	// This part will fail if those can't be loaded.
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(utils.Getenv("AWS_S3_BASE", ""))
		o.UsePathStyle = utils.Getenv("AWS_S3_USE_PATH_STYLE", "false") == "true"
	})

	return client
}
