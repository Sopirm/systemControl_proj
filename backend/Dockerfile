# Multi-stage Dockerfile for Go backend

# 1) Builder stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Enable CGO off for static binary and set target OS/ARCH
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Install git (for go mod) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Cache go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the backend source
COPY backend/ ./

# Build server and migrate utility
RUN go build -o server ./main.go \
 && go build -o migrate ./cmd/migrate/main.go

# 2) Runtime stage
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -H appuser

# Copy binary from builder
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate

# Environment defaults (can be overridden in docker-compose)
ENV SERVER_PORT=8080 \
    DB_HOST=db \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=systemcontrol \
    DB_SSLMODE=disable \
    JWT_SECRET=change_me \
    JWT_EXPIRATION_HRS=24

USER appuser
EXPOSE 8080

CMD ["/app/server"]
