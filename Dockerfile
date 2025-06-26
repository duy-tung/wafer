FROM golang:1.23 AS build

WORKDIR /src
COPY . .

# Build wafer binary with trimpath and stripped symbols
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /wafer ./cmd/wafer

# Final minimal image
FROM debian:bookworm-slim

# Install CA certificates and ensure update
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Add wafer binary
COPY --from=build /wafer /usr/bin/wafer

# Create storage directory for default output
RUN mkdir -p /app/storage

# Set working directory
WORKDIR /app

# Set default entrypoint so all subcommands are supported
ENTRYPOINT ["/usr/bin/wafer"]
