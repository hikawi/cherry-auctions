#!/bin/bash

SMTP_HOST=localhost \
  SMTP_PORT=1025 \
  SMTP_USER=test@example.com \
  SMTP_PASSWORD=test \
  AWS_ACCESS_KEY_ID=82d0d2b9a0cc2ba07bcb \
  AWS_SECRET_ACCESS_KEY=pX8uRMhyG26ENPoTUSZY9ziwxba0VrQKdtefOqcA \
  AWS_S3_BASE=http://localhost:9000 \
  AWS_S3_USE_PATH_STYLE=true \
  AWS_BUCKET_NAME=cherry-auctions \
  GIN_MODE=release \
  go test ./... -cover
