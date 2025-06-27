# Build stage - builds from source (for development)
FROM golang:1.24 AS build

WORKDIR /src

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build wafer binary with trimpath and stripped symbols
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /wafer ./cmd/wafer

# Release stage - copies pre-built binary (for GoReleaser)
FROM debian:bookworm-slim AS release

# Install CA certificates and ensure update
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Add wafer binary - will be from build stage OR copied by GoReleaser
COPY wafer /usr/bin/wafer

# Create storage directory for default output
RUN mkdir -p /app/storage

# Set working directory
WORKDIR /app

# Set default entrypoint so all subcommands are supported
ENTRYPOINT ["/usr/bin/wafer"]

# Default stage - builds from source (for docker build .)
FROM debian:bookworm-slim

# Install CA certificates and ensure update
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Add wafer binary from build stage
COPY --from=build /wafer /usr/bin/wafer

# Create storage directory for default output
RUN mkdir -p /app/storage

# Set working directory
WORKDIR /app

# Set default entrypoint so all subcommands are supported
ENTRYPOINT ["/usr/bin/wafer"]
