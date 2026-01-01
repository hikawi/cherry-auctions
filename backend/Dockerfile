FROM golang:1.25.4-alpine AS builder

# Some random ass things
# I don't even know what is needed
RUN apk add --no-cache vips-dev vips pkgconf gcc musl musl-dev build-base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
RUN go build -o ./app .

FROM alpine:3.22 AS prod

# Same random ass things
RUN apk add --no-cache vips pkgconf

WORKDIR /app
COPY --from=builder ./app/app ./app

EXPOSE 80
CMD ["./app"]

