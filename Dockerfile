# step 1: build executable binary
FROM golang:1.18-alpine AS builder
LABEL maintainer="Mikhael Kristian <kristianmikhael667@gmail.com>"
RUN apk --no-cache add git build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /qgames-user-service

# step 2: build a small image
FROM alpine:3.16.0
WORKDIR /app
COPY --from=builder /qgames-user-service .
COPY .env .

# Install nginx in the final image
RUN apk --no-cache add nginx

# Create the /run/nginx/ directory
RUN mkdir -p /run/nginx


EXPOSE 3001
CMD ["./qgames-user-service", "-migrate=migrate"]
