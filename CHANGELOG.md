# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Golden File Testing Framework**: Comprehensive regression testing with golden files
- **Enhanced Build System**: Sophisticated Makefile with colored output, emojis, and parallel execution
- **Performance Benchmarking**: Complete benchmark suite with baseline comparison and profiling
- **Comprehensive Documentation**: ROADMAP.md, BENCHMARK.md, CONTRIBUTING.md following mochi patterns
- **Development Tools**: Enhanced development setup with tool detection and environment checking
- **Docker containerization support** with multi-stage builds
- **GitHub Container Registry publishing workflow**
- **Docker Compose examples** and usage documentation
- **Makefile targets for Docker operations** (`docker-build`, `docker-test`, `docker-run`, `docker-push`, `docker-clean`)
- **Enhanced help system** with categorized targets and quick start examples
- **Benchmark targets** (`bench`, `bench-baseline`, `bench-compare`, `bench-profile`)
- **Golden file test targets** (`test-golden`, `update-golden`)
- **Development environment tools** (`dev-setup`, `dev-check`, `lint-install`)

### Enhanced
- **Makefile**: Complete rewrite with mochi's sophisticated patterns
  - Colored output with emojis for better UX
  - Parallel execution for build targets
  - Comprehensive help system with categorized targets
  - Tool detection with graceful fallbacks
  - Enhanced error handling and user feedback
- **Testing Strategy**:
  - Golden file testing framework for regression testing
  - Comprehensive benchmark suite for performance tracking
  - Enhanced unit and integration tests
  - Memory allocation and performance profiling
- **Documentation Structure**:
  - Enhanced README.md with comprehensive sections and better organization
  - Technical ROADMAP.md with detailed future planning
  - Performance BENCHMARK.md with detailed metrics and analysis
  - Professional CONTRIBUTING.md with clear guidelines
  - Improved project documentation following mochi standards

### Changed
- **Updated GoReleaser configuration** to support Docker image publishing
- **Enhanced installation and usage guides** with Docker instructions
- **Added Docker examples** to usage documentation
- **Restructured README.md** with better organization and comprehensive content
- **Enhanced SPEC.md** with more detailed technical specifications
- **Improved project structure** following mochi's proven patterns

## [0.1.0] - 2025-01-15

### Added
- Initial release of wafer CLI tool
- Core `ingest` command for processing text files
- Recursive directory traversal for `.txt` file discovery
- Text chunking with configurable word count (default: 300 words)
- Word boundary preservation in text chunking
- Ollama API integration for embedding generation
- Default support for `nomic-embed-text` model
- JSONL output format with comprehensive metadata
- UUID generation for each text chunk
- RFC3339 timestamp generation
- Structured logging with `log/slog`
- Comprehensive error handling and recovery
- Progress reporting for large directories
- Health check for Ollama API connectivity
- Cross-platform binary builds (Linux, macOS, Windows)
- Multiple architecture support (amd64, arm64)

### Features
- **CLI Interface**: Clean command-line interface with `--model`, `--output`, and `--chunk-size` flags
- **File Processing**: Handles UTF-8 encoded text files with proper error handling
- **Chunking Algorithm**: Smart text splitting that preserves word boundaries
- **API Integration**: Robust HTTP client with retry logic and exponential backoff
- **Output Format**: Structured JSONL with id, source_file, chunk_index, text, embedding, word_count, and created_at fields
- **Error Handling**: Graceful handling of file permission issues, API failures, and network problems
- **Logging**: Structured logging with different levels (DEBUG, INFO, WARN, ERROR)
- **Testing**: Comprehensive unit and integration test suite
- **Build System**: Professional Makefile with multiple targets
- **CI/CD**: GitHub Actions workflow for automated testing and releases
- **Documentation**: Complete README, technical specification, and usage guides

### Technical Details
- Built with Go 1.24
- Uses `github.com/alexflint/go-arg` for CLI argument parsing
- Uses `github.com/google/uuid` for UUID generation
- HTTP client with 30-second timeout and 3 retry attempts
- Supports relative and absolute directory paths
- Creates output directories automatically
- Append mode for incremental processing
- Memory-efficient processing (doesn't load entire files into memory)

### Dependencies
- `github.com/alexflint/go-arg` v1.5.1 - CLI argument parsing
- `github.com/google/uuid` v1.6.0 - UUID generation
- Go standard library for HTTP, JSON, file I/O, and logging

### Build and Release
- Cross-platform compilation support
- GoReleaser configuration for automated releases
- GitHub Actions CI/CD pipeline
- Static binary generation with embedded version information
- Support for Linux, macOS, and Windows on multiple architectures

[Unreleased]: https://github.com/duy-tung/wafer/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/duy-tung/wafer/releases/tag/v0.1.0
