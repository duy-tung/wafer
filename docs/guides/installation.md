# Installation Guide

This guide covers different ways to install and set up the wafer CLI tool.

## Prerequisites

### Ollama

Wafer requires [Ollama](https://ollama.ai/) to be running locally for embedding generation.

**Install Ollama:**

```bash
# Linux/macOS
curl -fsSL https://ollama.ai/install.sh | sh

# Windows (PowerShell)
iex (irm https://ollama.ai/install.ps1)
```

**Pull the default embedding model:**

```bash
ollama pull nomic-embed-text
```

**Start Ollama (if not already running):**

```bash
ollama serve
```

## Installation Methods

### Method 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/duy-tung/wafer/releases).

**Linux (x86_64):**
```bash
curl -L https://github.com/duy-tung/wafer/releases/latest/download/wafer-linux-amd64 -o wafer
chmod +x wafer
sudo mv wafer /usr/local/bin/
```

**Linux (ARM64):**
```bash
curl -L https://github.com/duy-tung/wafer/releases/latest/download/wafer-linux-arm64 -o wafer
chmod +x wafer
sudo mv wafer /usr/local/bin/
```

**macOS (Intel):**
```bash
curl -L https://github.com/duy-tung/wafer/releases/latest/download/wafer-darwin-amd64 -o wafer
chmod +x wafer
sudo mv wafer /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/duy-tung/wafer/releases/latest/download/wafer-darwin-arm64 -o wafer
chmod +x wafer
sudo mv wafer /usr/local/bin/
```

**Windows:**
```powershell
# Download wafer-windows-amd64.exe from releases page
# Add to PATH or place in a directory that's already in PATH
```

### Method 2: Build from Source

**Prerequisites:**
- Go 1.23 or later
- Make (optional, but recommended)
- Git

**Steps:**

```bash
# Clone the repository
git clone https://github.com/duy-tung/wafer.git
cd wafer

# Build using Make (recommended)
make build

# Or build manually
go build -o ~/bin/wafer cmd/wafer/main.go
```

### Method 3: Docker (Containerized)

**Prerequisites:**
- Docker installed and running

**Pull from GitHub Container Registry:**
```bash
docker pull ghcr.io/duy-tung/wafer:latest
```

**Build locally:**
```bash
git clone https://github.com/duy-tung/wafer.git
cd wafer
make docker-build
```

### Method 4: Go Install

If you have Go installed, you can install directly:

```bash
go install github.com/duy-tung/wafer/cmd/wafer@latest
```

## Verification

### Native Binary

Verify the installation by running:

```bash
wafer --version
```

You should see output similar to:
```
wafer v0.1.0 (abc1234, built 2024-01-15T10:30:45Z)
```

Test the help command:

```bash
wafer --help
```

### Docker

Verify the Docker image:

```bash
docker run --rm ghcr.io/duy-tung/wafer:latest --version
docker run --rm ghcr.io/duy-tung/wafer:latest --help
```

## Configuration

### Environment Variables

You can configure wafer using environment variables:

```bash
# Set custom Ollama host (default: http://localhost:11434)
export OLLAMA_HOST=http://your-ollama-server:11434

# Set log level (default: info)
export WAFER_LOG_LEVEL=debug
```

### Ollama Models

Install additional embedding models if needed:

```bash
# High-quality models
ollama pull nomic-embed-text    # Default, recommended
ollama pull mxbai-embed-large   # High accuracy

# Lightweight models
ollama pull all-minilm          # Fast, smaller
```

## Troubleshooting

### Common Issues

**1. "command not found: wafer"**
- Ensure the binary is in your PATH
- Check that you have execute permissions: `chmod +x wafer`

**2. "Ollama API is not accessible"**
- Verify Ollama is running: `ollama list`
- Check the Ollama service: `curl http://localhost:11434/api/tags`
- Ensure the correct host is configured

**3. "failed to get embedding"**
- Verify the model is installed: `ollama list`
- Pull the model if missing: `ollama pull nomic-embed-text`
- Check Ollama logs for errors

**4. Permission denied errors**
- Ensure you have read access to input directories
- Ensure you have write access to output directories
- Check file permissions: `ls -la`

### Getting Help

If you encounter issues:

1. Check the [troubleshooting section](../features/embedding.md#troubleshooting)
2. Review the [technical specification](../../SPEC.md)
3. Open an issue on [GitHub](https://github.com/duy-tung/wafer/issues)

## Next Steps

- Read the [Usage Guide](usage.md) to learn how to use wafer
- Explore [Embedding Features](../features/embedding.md) for advanced usage
- Check out the [examples](../../tests/fixtures/) for sample inputs
