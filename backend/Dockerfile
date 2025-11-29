FROM golang:1.25.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o ./app ./cmd/server/main.go

FROM alpine:3.22 AS prod

WORKDIR /app
COPY --from=builder ./app/app ./app

EXPOSE 80
CMD ["./app"]

