package test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/utils"
)

func (h *TestHandler) GetTest(g *gin.Context) {
	ctx := context.Background()
	res, err := h.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		utils.Log(gin.H{"error": err})
		g.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		return
	}

	names := make([]string, 0)
	for _, bucket := range res.Buckets {
		names = append(names, *bucket.Name)
	}

	g.JSON(200, names)
}
