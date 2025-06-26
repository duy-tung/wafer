# Default target
.DEFAULT_GOAL := help
.PHONY: build test test-golden test-unit test-integration lint fmt clean install release help bench dev-setup update-golden

# Project metadata
APP_NAME := wafer
BIN_DIR := $(HOME)/bin
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Source paths
MAIN_SRC := cmd/wafer/main.go

# Build flags
GO := go
GOFLAGS := -trimpath -mod=readonly
LDFLAGS := "-s -w -X 'main.version=$(VERSION)' -X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'"

# Colors and emojis for output (following mochi pattern)
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
MAGENTA := \033[35m
RESET := \033[0m

# Emojis for better UX
ROCKET := 🚀
GEAR := ⚙️
TEST_TUBE := 🧪
MAGNIFYING_GLASS := 🔍
BROOM := 🧹
PACKAGE := 📦
SPARKLES := ✨
CHECKMARK := ✅
CROSS := ❌
HOURGLASS := ⏳

# --------------------------
# Build Targets
# --------------------------

build: ## Build binary for current platform
	@echo "$(CYAN)$(GEAR) Building $(APP_NAME)...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_SRC)
	@echo "$(GREEN)$(CHECKMARK) Built: $(BIN_DIR)/$(APP_NAME) (v$(VERSION))$(RESET)"

build-all: ## Build for all platforms (parallel)
	@echo "$(CYAN)$(GEAR) Building for all platforms...$(RESET)"
	@mkdir -p dist
	@$(MAKE) -j4 build-linux build-darwin build-windows
	@echo "$(GREEN)$(CHECKMARK) All platforms built$(RESET)"

build-linux: ## Build for Linux
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-linux-amd64 $(MAIN_SRC)

build-darwin: ## Build for macOS
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-darwin-amd64 $(MAIN_SRC)

build-windows: ## Build for Windows
	@GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-windows-amd64.exe $(MAIN_SRC)

# --------------------------
# Testing
# --------------------------

test: ## Run all tests with coverage
	@echo "$(CYAN)$(TEST_TUBE) Running all tests...$(RESET)"
	@$(GO) test ./... -v -race -coverprofile=coverage.out
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)$(CHECKMARK) Tests completed. Coverage report: coverage.html$(RESET)"

test-unit: ## Run only unit tests
	@echo "$(CYAN)$(TEST_TUBE) Running unit tests...$(RESET)"
	@$(GO) test ./tests -run TestChunker -v
	@$(GO) test ./tests -run TestEmbedder -v
	@$(GO) test ./tests -run TestWriter -v
	@$(GO) test ./tests -run TestConfig -v
	@echo "$(GREEN)$(CHECKMARK) Unit tests completed$(RESET)"

test-integration: ## Run only integration tests
	@echo "$(CYAN)$(TEST_TUBE) Running integration tests...$(RESET)"
	@$(GO) test ./tests -run TestEndToEnd -v
	@echo "$(GREEN)$(CHECKMARK) Integration tests completed$(RESET)"

test-golden: ## Run golden file tests
	@echo "$(CYAN)$(TEST_TUBE) Running golden file tests...$(RESET)"
	@$(GO) test ./tests -run TestGoldenFiles -v
	@echo "$(GREEN)$(CHECKMARK) Golden file tests completed$(RESET)"

update-golden: ## Update golden files
	@echo "$(YELLOW)$(HOURGLASS) Updating golden files...$(RESET)"
	@$(GO) test ./tests -run TestGoldenFiles -update
	@echo "$(GREEN)$(CHECKMARK) Golden files updated$(RESET)"

test-fast: ## Run tests without race detection (faster)
	@echo "$(CYAN)$(TEST_TUBE) Running fast tests...$(RESET)"
	@$(GO) test ./... -v
	@echo "$(GREEN)$(CHECKMARK) Fast tests completed$(RESET)"

# --------------------------
# Code Quality
# --------------------------

fmt: ## Format source files
	@echo "$(CYAN)$(BROOM) Formatting source...$(RESET)"
	@$(GO) fmt ./...
	@echo "$(GREEN)$(CHECKMARK) Code formatted$(RESET)"

lint: ## Run static analysis
	@echo "$(CYAN)$(MAGNIFYING_GLASS) Running static analysis...$(RESET)"
	@if command -v golangci-lint > /dev/null 2>&1; then \
		echo "$(CYAN)$(PACKAGE) Using golangci-lint$(RESET)"; \
		golangci-lint run ./... --timeout=5m; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found, falling back to go vet$(RESET)"; \
		$(GO) vet ./...; \
	fi
	@echo "$(GREEN)$(CHECKMARK) Linting completed$(RESET)"

lint-install: ## Install golangci-lint
	@echo "$(CYAN)$(PACKAGE) Installing golangci-lint...$(RESET)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2
	@echo "$(GREEN)$(CHECKMARK) golangci-lint installed$(RESET)"

# --------------------------
# Benchmarking
# --------------------------

bench: ## Run benchmarks
	@echo "$(CYAN)$(HOURGLASS) Running benchmarks...$(RESET)"
	@$(GO) test -bench=. -benchmem ./... | tee bench/results.txt
	@echo "$(GREEN)$(CHECKMARK) Benchmarks completed. Results: bench/results.txt$(RESET)"

bench-compare: ## Compare benchmarks with previous run
	@echo "$(CYAN)$(HOURGLASS) Comparing benchmarks...$(RESET)"
	@if [ -f bench/baseline.txt ]; then \
		benchcmp bench/baseline.txt bench/results.txt; \
	else \
		echo "$(YELLOW)⚠️  No baseline found. Run 'make bench-baseline' first$(RESET)"; \
	fi

bench-baseline: ## Set current benchmark as baseline
	@echo "$(CYAN)$(HOURGLASS) Setting benchmark baseline...$(RESET)"
	@mkdir -p bench
	@$(GO) test -bench=. -benchmem ./... > bench/baseline.txt
	@echo "$(GREEN)$(CHECKMARK) Baseline set$(RESET)"

bench-profile: ## Run benchmarks with CPU profiling
	@echo "$(CYAN)$(HOURGLASS) Running benchmarks with profiling...$(RESET)"
	@mkdir -p bench
	@$(GO) test -bench=. -benchmem -cpuprofile=bench/cpu.prof ./...
	@echo "$(GREEN)$(CHECKMARK) Profiling completed. Profile: bench/cpu.prof$(RESET)"

# --------------------------
# Dependencies
# --------------------------

deps: ## Download and tidy dependencies
	@echo "$(CYAN)$(PACKAGE) Managing dependencies...$(RESET)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(GREEN)$(CHECKMARK) Dependencies updated$(RESET)"

# --------------------------
# Installation
# --------------------------

install: build ## Install binary to $GOPATH/bin
	@echo "$(GREEN)✅ $(APP_NAME) installed to $(BIN_DIR)$(RESET)"

# --------------------------
# Release
# --------------------------

release: ## Build release binaries for multiple platforms
	@echo "$(CYAN)🚀 Building release binaries...$(RESET)"
	@mkdir -p dist
	
	# Linux amd64
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-linux-amd64 $(MAIN_SRC)
	
	# Linux arm64
	@GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-linux-arm64 $(MAIN_SRC)
	
	# macOS amd64
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-darwin-amd64 $(MAIN_SRC)
	
	# macOS arm64
	@GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-darwin-arm64 $(MAIN_SRC)
	
	# Windows amd64
	@GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags=$(LDFLAGS) -o dist/$(APP_NAME)-windows-amd64.exe $(MAIN_SRC)
	
	@echo "$(GREEN)✅ Release binaries built in dist/$(RESET)"
	@ls -la dist/

# --------------------------
# Cleanup
# --------------------------

clean: ## Clean built binaries and artifacts
	@echo "$(CYAN)🧽 Cleaning artifacts...$(RESET)"
	@rm -f $(BIN_DIR)/$(APP_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✅ Cleanup completed$(RESET)"

clean-all: clean docker-clean ## Clean everything including Docker images
	@echo "$(GREEN)✅ Complete cleanup finished$(RESET)"

# --------------------------
# Docker
# --------------------------

docker-build: ## Build Docker image
	@echo "$(CYAN)🐳 Building Docker image...$(RESET)"
	@docker build -t wafer:latest .
	@docker build -t wafer:$(VERSION) .
	@echo "$(GREEN)✅ Docker image built: wafer:latest, wafer:$(VERSION)$(RESET)"

docker-run: docker-build ## Build and run Docker container with test fixtures
	@echo "$(CYAN)🐳 Running Docker container with test fixtures...$(RESET)"
	@docker run --rm -v $(PWD)/tests/fixtures:/data -v $(PWD)/storage:/app/storage wafer:latest ingest /data --chunk-size=50
	@echo "$(GREEN)✅ Docker run completed. Check storage/ for output$(RESET)"

docker-test: docker-build ## Test Docker image
	@echo "$(CYAN)🐳 Testing Docker image...$(RESET)"
	@docker run --rm wafer:latest --version
	@docker run --rm wafer:latest --help
	@echo "$(GREEN)✅ Docker image tests passed$(RESET)"

docker-push: docker-build ## Push Docker image to registry (requires login)
	@echo "$(CYAN)🐳 Pushing Docker image...$(RESET)"
	@docker tag wafer:latest ghcr.io/duy-tung/wafer:latest
	@docker tag wafer:$(VERSION) ghcr.io/duy-tung/wafer:$(VERSION)
	@docker push ghcr.io/duy-tung/wafer:latest
	@docker push ghcr.io/duy-tung/wafer:$(VERSION)
	@echo "$(GREEN)✅ Docker images pushed$(RESET)"

docker-clean: ## Clean Docker images
	@echo "$(CYAN)🐳 Cleaning Docker images...$(RESET)"
	@docker rmi wafer:latest wafer:$(VERSION) 2>/dev/null || true
	@docker rmi ghcr.io/duy-tung/wafer:latest ghcr.io/duy-tung/wafer:$(VERSION) 2>/dev/null || true
	@echo "$(GREEN)✅ Docker images cleaned$(RESET)"

# --------------------------
# Development
# --------------------------

dev-setup: ## Set up development environment
	@echo "$(CYAN)$(GEAR) Setting up development environment...$(RESET)"
	@$(GO) mod download
	@mkdir -p bench golden tests/golden
	@if ! command -v golangci-lint > /dev/null 2>&1; then \
		echo "$(YELLOW)$(PACKAGE) Installing golangci-lint...$(RESET)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	fi
	@if ! command -v benchcmp > /dev/null 2>&1; then \
		echo "$(YELLOW)$(PACKAGE) Installing benchcmp...$(RESET)"; \
		$(GO) install golang.org/x/tools/cmd/benchcmp@latest; \
	fi
	@echo "$(GREEN)$(CHECKMARK) Development environment ready$(RESET)"

dev-check: ## Check development environment
	@echo "$(CYAN)$(MAGNIFYING_GLASS) Checking development environment...$(RESET)"
	@echo "Go version: $$($(GO) version)"
	@echo "golangci-lint: $$(if command -v golangci-lint > /dev/null 2>&1; then echo '$(GREEN)$(CHECKMARK) installed$(RESET)'; else echo '$(RED)$(CROSS) not found$(RESET)'; fi)"
	@echo "benchcmp: $$(if command -v benchcmp > /dev/null 2>&1; then echo '$(GREEN)$(CHECKMARK) installed$(RESET)'; else echo '$(RED)$(CROSS) not found$(RESET)'; fi)"
	@echo "Docker: $$(if command -v docker > /dev/null 2>&1; then echo '$(GREEN)$(CHECKMARK) installed$(RESET)'; else echo '$(RED)$(CROSS) not found$(RESET)'; fi)"

run-example: build ## Build and run example with test fixtures
	@echo "$(CYAN)🏃 Running example with test fixtures...$(RESET)"
	@mkdir -p storage
	@$(BIN_DIR)/$(APP_NAME) ingest tests/fixtures --output storage/example_vectors.jsonl --chunk-size 50
	@echo "$(GREEN)✅ Example completed. Check storage/example_vectors.jsonl$(RESET)"

# --------------------------
# Help
# --------------------------

help: ## Show help message
	@echo ""
	@echo "$(CYAN)$(PACKAGE) Wafer CLI Tool Makefile$(RESET)"
	@echo "$(CYAN)==============================$(RESET)"
	@echo ""
	@echo "$(BLUE)$(GEAR) Build Targets:$(RESET)"
	@grep -E '^(build|release).*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BLUE)$(TEST_TUBE) Testing Targets:$(RESET)"
	@grep -E '^(test|update-golden).*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BLUE)$(MAGNIFYING_GLASS) Quality Targets:$(RESET)"
	@grep -E '^(lint|fmt).*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BLUE)$(HOURGLASS) Benchmark Targets:$(RESET)"
	@grep -E '^bench.*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BLUE)🐳 Docker Targets:$(RESET)"
	@grep -E '^docker.*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BLUE)$(GEAR) Development Targets:$(RESET)"
	@grep -E '^(dev-|deps|clean).*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)$(SPARKLES) Quick Start Examples:$(RESET)"
	@echo "  $(GREEN)make dev-setup$(RESET)      # Set up development environment"
	@echo "  $(GREEN)make build$(RESET)          # Build the binary"
	@echo "  $(GREEN)make test$(RESET)           # Run all tests"
	@echo "  $(GREEN)make test-golden$(RESET)    # Run golden file tests"
	@echo "  $(GREEN)make bench$(RESET)          # Run benchmarks"
	@echo "  $(GREEN)make docker-build$(RESET)   # Build Docker image"
	@echo "  $(GREEN)make release$(RESET)        # Build for all platforms"
	@echo ""
