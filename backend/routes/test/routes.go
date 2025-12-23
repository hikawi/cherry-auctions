package test

import "github.com/aws/aws-sdk-go-v2/service/s3"

type TestHandler struct {
	S3 *s3.Client
}
