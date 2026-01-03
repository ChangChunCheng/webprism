# WEBPRISM Dockerfile
# Multi-stage build for optimized image size

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies (protoc and build tools)
RUN apk add --no-cache git make protobuf-dev

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Install protoc-gen tools to GOPATH/bin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Copy source code
COPY . .

# Generate protobuf code
RUN make proto

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webprism-server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webprism-mcp cmd/mcp/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webprism cmd/cli/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/webprism-server .
COPY --from=builder /build/webprism-mcp .
COPY --from=builder /build/webprism .

# Copy configuration
COPY deployments/config.yaml /app/deployments/config.yaml

# Create non-root user
RUN addgroup -g 1000 webprism && \
    adduser -D -u 1000 -G webprism webprism && \
    chown -R webprism:webprism /app

USER webprism

# Expose ports
EXPOSE 8080 9090

# Default command
CMD ["./webprism-server"]
