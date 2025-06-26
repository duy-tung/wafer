# Wafer Technical Specification

## Overview

Wafer is a production-ready Go CLI tool designed to process text files and generate embeddings using Ollama. It follows modern CLI design patterns and provides comprehensive error handling, logging, and progress reporting.

## Architecture

### Project Structure

```
wafer/
├── cmd/wafer/main.go           # CLI entrypoint and argument parsing
├── internal/
│   ├── ingest/
│   │   ├── processor.go        # Core orchestration logic
│   │   ├── chunker.go         # Text chunking implementation
│   │   ├── embedder.go        # Ollama API client
│   │   └── writer.go          # JSONL output handling
│   └── config/
│       └── config.go          # Configuration management
├── storage/                    # Default output directory
├── tests/                      # Test suite
└── docs/                       # Documentation
```

### Core Components

#### 1. CLI Interface (`cmd/wafer/main.go`)

**Command Structure:**
```bash
wafer ingest <directory_path> [--model=MODEL_NAME] [--output=OUTPUT_PATH] [--chunk-size=WORDS]
```

**Argument Parsing:**
- Uses `github.com/alexflint/go-arg` for structured argument parsing
- Supports subcommands with typed arguments
- Provides built-in help and version information

**Configuration:**
- Default model: `nomic-embed-text`
- Default output: `storage/vectors.jsonl`
- Default chunk size: `300` words

#### 2. Configuration Management (`internal/config/config.go`)

**Config Structure:**
```go
type Config struct {
    Directory string // Directory to process
    Model     string // Ollama model name
    Output    string // Output file path
    ChunkSize int    // Chunk size in words
}
```

**Validation:**
- Directory existence and accessibility
- Positive chunk size validation
- Output directory creation
- Model name validation

#### 3. Text Processing (`internal/ingest/chunker.go`)

**Chunking Algorithm:**
- Splits text into approximately N-word chunks (configurable)
- Preserves word boundaries using `bufio.ScanWords`
- Handles UTF-8 encoding properly
- Filters out non-alphanumeric tokens
- Maintains chunk indices for ordering

**Edge Cases:**
- Files smaller than chunk size → single chunk
- Empty files → no chunks
- Files with only whitespace → no chunks
- Unicode characters → properly handled

#### 4. Ollama Integration (`internal/ingest/embedder.go`)

**API Client:**
- HTTP client with 30-second timeout
- Exponential backoff retry mechanism (3 attempts)
- Context-aware cancellation support
- Health check functionality

**Request Format:**
```json
{
  "model": "nomic-embed-text",
  "prompt": "text chunk content"
}
```

**Response Format:**
```json
{
  "embedding": [0.1234, -0.5678, 0.9012, ...]
}
```

**Error Handling:**
- Network connectivity issues
- API rate limiting
- Invalid model names
- Malformed responses

#### 5. Output Generation (`internal/ingest/writer.go`)

**JSONL Schema:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "source_file": "relative/path/to/file.txt",
  "chunk_index": 0,
  "text": "actual chunk text content here...",
  "embedding": [0.1234, -0.5678, 0.9012, ...],
  "word_count": 299,
  "created_at": "2024-01-15T10:30:45Z"
}
```

**Field Specifications:**
- `id`: UUID v4 generated for each chunk
- `source_file`: Relative path from input directory
- `chunk_index`: Zero-based sequential index
- `text`: Original chunk text content
- `embedding`: Float64 array from Ollama
- `word_count`: Actual word count in chunk
- `created_at`: RFC3339 UTC timestamp

**File Operations:**
- Append mode for incremental processing
- Automatic directory creation
- Proper file handle management
- Atomic write operations

#### 6. Process Orchestration (`internal/ingest/processor.go`)

**Workflow:**
1. Configuration validation
2. Ollama API health check
3. Recursive file discovery
4. File processing with progress reporting
5. Statistics collection and reporting

**File Discovery:**
- Uses `filepath.WalkDir` for recursive traversal
- Filters for `.txt` extension (case-insensitive)
- Handles permission errors gracefully
- Skips directories and non-text files

**Error Handling:**
- Individual file failures don't stop processing
- Comprehensive error logging with context
- Final statistics include error counts
- Appropriate exit codes

## API Specifications

### Ollama API Integration

**Base URL:** `http://localhost:11434` (configurable)

**Endpoints Used:**
- `POST /api/embeddings` - Generate embeddings
- `GET /api/tags` - Health check

**Request Headers:**
- `Content-Type: application/json`

**Timeout Configuration:**
- HTTP client timeout: 30 seconds
- Retry attempts: 3
- Backoff strategy: Exponential (1s, 2s, 4s)

### Output Format Specification

**File Format:** JSONL (JSON Lines)
**Encoding:** UTF-8
**Line Termination:** Unix-style (`\n`)

**Required Fields:**
- All fields are required and must be non-empty/non-null
- `embedding` array must contain at least one element
- `word_count` must be positive integer
- `chunk_index` must be non-negative integer

## Performance Characteristics

### Memory Usage
- Base application: ~10MB
- Per chunk processing: ~1KB
- Ollama model memory: Variable (model-dependent)

### Throughput
- File I/O: Limited by disk speed
- Network: Limited by Ollama API response time
- Typical: 100-500 chunks/minute

### Scalability
- Processes files sequentially (single-threaded)
- Memory usage scales with chunk size, not file count
- Output file grows linearly with input

## Error Handling Strategy

### Error Categories

1. **Configuration Errors**
   - Invalid directory paths
   - Permission issues
   - Invalid parameters

2. **Runtime Errors**
   - File I/O failures
   - Network connectivity issues
   - API failures

3. **Data Errors**
   - Malformed text files
   - Empty files
   - Encoding issues

### Error Response

**Logging Levels:**
- `ERROR`: Critical failures that stop processing
- `WARN`: Non-critical issues (skipped files)
- `INFO`: Normal operation progress
- `DEBUG`: Detailed operation information

**Exit Codes:**
- `0`: Success
- `1`: Configuration or critical error
- `2`: Partial failure (some files processed)

## Testing Strategy

### Unit Tests
- Individual component testing
- Mock external dependencies
- Edge case validation
- Error condition testing

### Integration Tests
- End-to-end workflow testing
- Mock Ollama server
- Golden file comparisons
- Multi-file processing

### Test Coverage
- Target: >90% code coverage
- Critical paths: 100% coverage
- Error handling: Comprehensive testing

## Security Considerations

### Input Validation
- Path traversal prevention
- File size limits (configurable)
- Content sanitization

### Network Security
- HTTPS support for Ollama API
- Timeout enforcement
- Request size limits

### Output Security
- Safe file creation
- Directory traversal prevention
- Atomic write operations

## Dependencies

### Core Dependencies
- `github.com/alexflint/go-arg` - CLI argument parsing
- `github.com/google/uuid` - UUID generation
- Standard library: `net/http`, `log/slog`, `encoding/json`

### Development Dependencies
- `golangci-lint` - Static analysis
- `go test` - Testing framework
- `make` - Build automation

## Build and Release

### Build Process
- Go modules for dependency management
- Cross-platform compilation support
- Version embedding via ldflags
- Static binary generation
- Multi-stage Docker builds with build caching

### Release Process
- GitHub Actions CI/CD
- GoReleaser for multi-platform builds
- Automated testing on multiple OS
- Docker image publishing to GitHub Container Registry
- Semantic versioning

### Supported Platforms
- Linux: amd64, arm64
- macOS: amd64, arm64
- Windows: amd64
- Docker: linux/amd64, linux/arm64

### Docker Support
- Multi-stage Dockerfile for minimal production images
- Base image: debian:bookworm-slim with CA certificates
- Build caching for faster subsequent builds
- Volume mounting for input/output directories
- Network host mode for Ollama connectivity

## Docker Containerization

### Container Architecture
- **Base Image**: debian:bookworm-slim (minimal, secure)
- **Build Stage**: golang:1.24 (multi-stage build for smaller final image)
- **Binary Size**: ~6MB statically linked binary
- **Final Image Size**: ~80MB (including CA certificates)

### Container Features
- **Security**: Non-root execution, minimal attack surface
- **Networking**: Host network mode for Ollama connectivity
- **Storage**: Volume mounting for input/output directories
- **Caching**: Build cache optimization for faster rebuilds

### Usage Patterns
```bash
# Basic usage with volume mounts
docker run --rm \
  -v ./input:/data \
  -v ./output:/app/storage \
  --network host \
  ghcr.io/duy-tung/wafer:latest \
  ingest /data

# Docker Compose deployment
version: '3.8'
services:
  wafer:
    image: ghcr.io/duy-tung/wafer:latest
    volumes:
      - ./documents:/data:ro
      - ./output:/app/storage
    command: ingest /data --chunk-size=300
    network_mode: host
```

### Registry Information
- **Registry**: GitHub Container Registry (ghcr.io)
- **Repository**: ghcr.io/duy-tung/wafer
- **Tags**: latest, version-specific (e.g., v0.1.0)
- **Platforms**: linux/amd64, linux/arm64

## Configuration Options

### Command Line Flags
- `--model`: Ollama model name
- `--output`: Output file path
- `--chunk-size`: Words per chunk

### Environment Variables
- `OLLAMA_HOST`: Ollama server URL
- `WAFER_LOG_LEVEL`: Logging level

### Configuration Files
- Not currently supported
- Future enhancement possibility

## Monitoring and Observability

### Logging
- Structured logging with `log/slog`
- Configurable log levels
- Context-aware log messages
- Performance metrics logging

### Metrics
- Files processed count
- Chunks generated count
- Error counts by type
- Processing duration
- API response times

### Health Checks
- Ollama API connectivity
- Output directory writability
- Input directory accessibility
