FROM golang:1.25.4-alpine AS builder

# Some random ass things
RUN apk add --no-cache vips

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o ./app .

FROM alpine:3.22 AS prod

# Same random ass things
RUN apk add --no-cache vips

WORKDIR /app
COPY --from=builder ./app/app ./app

EXPOSE 80
CMD ["./app"]

