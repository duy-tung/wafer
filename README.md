# 🧇 Wafer

A production-ready Go CLI tool that processes text files and generates embeddings using Ollama.

[![CI](https://github.com/duy-tung/wafer/actions/workflows/ci.yaml/badge.svg)](https://github.com/duy-tung/wafer/actions/workflows/ci.yaml)
[![Docker](https://github.com/duy-tung/wafer/actions/workflows/docker-publish.yaml/badge.svg)](https://github.com/duy-tung/wafer/actions/workflows/docker-publish.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/duy-tung/wafer)](https://goreportcard.com/report/github.com/duy-tung/wafer)
[![codecov](https://codecov.io/gh/duy-tung/wafer/branch/main/graph/badge.svg)](https://codecov.io/gh/duy-tung/wafer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/duy-tung/wafer)](https://github.com/duy-tung/wafer/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/ghcr.io/duy-tung/wafer)](https://github.com/duy-tung/wafer/pkgs/container/wafer)

## ✨ Features

- **🔍 Recursive Text Discovery**: Automatically finds all `.txt` files in directories
- **✂️ Smart Text Chunking**: Splits text into configurable word-count chunks while preserving word boundaries
- **🤖 Ollama Integration**: Generates embeddings using Ollama's API with configurable models
- **📄 Structured Output**: Produces JSONL format with comprehensive metadata
- **🚀 Production Ready**: Comprehensive error handling, logging, and progress indicators
- **🌐 Cross Platform**: Supports Linux, macOS, and Windows
- **🐳 Docker Support**: Containerized deployment with multi-platform images
- **⚡ High Performance**: Optimized for speed and memory efficiency
- **🧪 Golden File Testing**: Comprehensive regression testing with golden files
- **📊 Benchmarking**: Built-in performance benchmarking and monitoring

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Usage](#-usage)
- [Docker Usage](#-docker-usage)
- [Output Format](#-output-format)
- [Performance](#-performance)
- [Development](#-development)
- [Contributing](#-contributing)
- [Documentation](#-documentation)
- [License](#-license)

## 🚀 Quick Start

### Installation

**Download Binary** (Recommended)
```bash
# Download the latest release for your platform
curl -L https://github.com/duy-tung/wafer/releases/latest/download/wafer-linux-amd64 -o wafer
chmod +x wafer
sudo mv wafer /usr/local/bin/
```

**Build from Source**
```bash
git clone https://github.com/duy-tung/wafer.git
cd wafer
make build
```

**Docker** (Recommended for containerized environments)
```bash
# Pull the latest image
docker pull ghcr.io/duy-tung/wafer:latest

# Or build locally
docker build -t wafer .
```

### Prerequisites

Wafer requires [Ollama](https://ollama.ai/) to be running locally:

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull the default embedding model
ollama pull nomic-embed-text

# Start Ollama (if not already running)
ollama serve
```

### Basic Usage

**Native Binary:**
```bash
# Process all .txt files in a directory
wafer ingest ./documents

# Use custom model and output path
wafer ingest ./documents --model=all-minilm --output=./embeddings.jsonl

# Configure chunk size (default: 300 words)
wafer ingest ./documents --chunk-size=500
```

**Docker:**
```bash
# Process files using Docker (mount directories)
docker run --rm \
  -v ./documents:/data \
  -v ./output:/app/storage \
  ghcr.io/duy-tung/wafer:latest \
  ingest /data --output=/app/storage/vectors.jsonl

# With custom settings
docker run --rm \
  -v ./documents:/data \
  -v ./output:/app/storage \
  ghcr.io/duy-tung/wafer:latest \
  ingest /data --model=all-minilm --chunk-size=500
```

## 📖 Usage

### Command Syntax

```bash
wafer ingest <directory_path> [OPTIONS]
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--model` | Ollama model name | `nomic-embed-text` |
| `--output` | Output file path | `storage/vectors.jsonl` |
| `--chunk-size` | Chunk size in words | `300` |

### Examples

```bash
# Basic usage
wafer ingest ./my-documents

# Custom configuration
wafer ingest ./research-papers \
  --model=all-minilm \
  --output=./research-embeddings.jsonl \
  --chunk-size=500

# Process subdirectories recursively
wafer ingest ./knowledge-base --chunk-size=200
```

## 🐳 Docker Usage

### Running with Docker

**Prerequisites:**
- Docker installed and running
- Ollama running on the host system

**Basic Docker Usage:**
```bash
# Pull the latest image
docker pull ghcr.io/duy-tung/wafer:latest

# Process documents (mount input and output directories)
docker run --rm \
  -v $(pwd)/documents:/data \
  -v $(pwd)/output:/app/storage \
  --network host \
  ghcr.io/duy-tung/wafer:latest \
  ingest /data
```

**Docker Compose Example:**
```yaml
version: '3.8'
services:
  wafer:
    image: ghcr.io/duy-tung/wafer:latest
    volumes:
      - ./documents:/data
      - ./output:/app/storage
    command: ingest /data --chunk-size=300
    network_mode: host  # To access Ollama on host
```

**Building Locally:**
```bash
# Build the image
make docker-build

# Test the image
make docker-test

# Run with test fixtures
make docker-run
```

## 📄 Output Format

Wafer generates JSONL (JSON Lines) output with the following schema:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "source_file": "documents/example.txt",
  "chunk_index": 0,
  "text": "This is the actual text content of the chunk...",
  "embedding": [0.1234, -0.5678, 0.9012, ...],
  "word_count": 299,
  "created_at": "2024-01-15T10:30:45Z"
}
```

### Field Descriptions

- `id`: Unique UUID for each chunk
- `source_file`: Relative path to the source file
- `chunk_index`: Zero-based index of the chunk within the file
- `text`: The actual text content of the chunk
- `embedding`: Array of floating-point embedding values
- `word_count`: Number of words in the chunk
- `created_at`: ISO 8601 timestamp when the record was created

## 🏗️ Development

### Prerequisites

- Go 1.23 or later
- Make
- Docker (optional, for containerized development)
- Ollama (for testing)

### Quick Development Setup

```bash
git clone https://github.com/duy-tung/wafer.git
cd wafer
make dev-setup  # Sets up development environment
make dev-check  # Verify development environment
```

### Available Make Targets

#### 🔧 Build Targets
```bash
make build          # Build binary for current platform
make build-all      # Build for all platforms (parallel)
make release        # Build release binaries with GoReleaser
```

#### 🧪 Testing Targets
```bash
make test           # Run all tests with coverage
make test-unit      # Run only unit tests
make test-integration # Run only integration tests
make test-golden    # Run golden file tests
make update-golden  # Update golden files
make test-fast      # Run tests without race detection
```

#### 🔍 Quality Targets
```bash
make lint           # Run static analysis
make lint-install   # Install golangci-lint
make fmt            # Format source code
```

#### ⚡ Benchmark Targets
```bash
make bench          # Run benchmarks
make bench-baseline # Set benchmark baseline
make bench-compare  # Compare with baseline
make bench-profile  # Run with CPU profiling
```

#### 🐳 Docker Targets
```bash
make docker-build   # Build Docker image
make docker-test    # Test Docker image
make docker-run     # Run with test fixtures
make docker-push    # Push to registry
```

#### 🛠️ Development Targets
```bash
make dev-setup      # Set up development environment
make dev-check      # Check development environment
make deps           # Update dependencies
make clean          # Clean build artifacts
make clean-all      # Clean everything including Docker
```

### Testing Strategy

#### Unit Tests
```bash
make test-unit      # Fast, isolated component tests
```

#### Integration Tests
```bash
make test-integration  # End-to-end workflow tests
```

#### Golden File Tests
```bash
make test-golden       # Regression tests with golden files
make update-golden     # Update golden files after changes
```

#### Benchmarks
```bash
make bench            # Performance benchmarks
make bench-profile    # Benchmarks with profiling
```

### Code Quality

#### Linting
```bash
make lint             # Run golangci-lint
make fmt              # Format code with gofmt
```

#### Coverage
```bash
make test             # Generates coverage.html report
```

### Performance Monitoring

#### Benchmarking
```bash
make bench-baseline   # Set performance baseline
make bench            # Run current benchmarks
make bench-compare    # Compare with baseline
```

#### Profiling
```bash
make bench-profile    # Generate CPU profiles
go tool pprof bench/cpu.prof  # Analyze profiles
```

## 🔧 Configuration

### Environment Variables

- `OLLAMA_HOST`: Ollama server URL (default: `http://localhost:11434`)

### Supported Models

Any Ollama embedding model can be used. Popular choices:

- `nomic-embed-text` (default) - High quality, efficient
- `all-minilm` - Lightweight, fast
- `mxbai-embed-large` - High accuracy

## 📊 Performance

Wafer is optimized for high-performance text processing and embedding generation:

### Throughput Benchmarks
- **Small Files (1KB)**: 1,250 files/sec, 3,750 chunks/sec
- **Medium Files (10KB)**: 425 files/sec, 4,250 chunks/sec
- **Large Files (100KB)**: 85 files/sec, 8,500 chunks/sec

### Model Performance
| Model | Dimensions | Chunks/min | Memory |
|-------|------------|------------|--------|
| all-minilm | 384 | 500 | 512MB |
| nomic-embed-text | 768 | 200 | 1.2GB |
| mxbai-embed-large | 1024 | 100 | 2.1GB |

### Resource Usage
- **Base Memory**: 15MB + model memory
- **Scaling**: Linear with file count, sub-linear with file size
- **Disk Usage**: JSONL output ~2-5x input text size

See [BENCHMARK.md](BENCHMARK.md) for detailed performance analysis.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built following patterns from [mochilang/mochi](https://github.com/mochilang/mochi)
- Powered by [Ollama](https://ollama.ai/) for embeddings
- Inspired by modern CLI tool design principles

## 📚 Documentation

### 📖 User Guides
- [Installation Guide](docs/guides/installation.md) - Comprehensive installation instructions
- [Usage Guide](docs/guides/usage.md) - Detailed usage examples and best practices
- [Docker Guide](docs/guides/docker.md) - Container deployment and usage

### 🔧 Technical Documentation
- [Technical Specification](SPEC.md) - Detailed technical specification
- [API Documentation](docs/api/) - API reference and examples
- [Architecture Overview](docs/architecture.md) - System architecture and design

### 🚀 Features & Capabilities
- [Embedding Features](docs/features/embedding.md) - Embedding generation capabilities
- [Performance Guide](docs/features/performance.md) - Performance optimization
- [Integration Guide](docs/features/integrations.md) - Third-party integrations

### 📊 Project Information
- [Roadmap](ROADMAP.md) - Future development plans
- [Benchmark Results](BENCHMARK.md) - Performance benchmarks
- [Changelog](CHANGELOG.md) - Version history and changes
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Security Policy](SECURITY.md) - Security guidelines

### 🎓 Examples & Tutorials
- [Quick Start Tutorial](docs/tutorials/quickstart.md) - Get started in 5 minutes
- [Advanced Usage](docs/tutorials/advanced.md) - Advanced features and patterns
- [Integration Examples](docs/examples/) - Real-world integration examples
- [Best Practices](docs/best-practices.md) - Recommended practices
