package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"wafer/golden"
)

func TestGoldenFiles(t *testing.T) {
	// Create mock Ollama server for golden tests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embeddings":
			// Return deterministic mock embedding for golden tests
			response := map[string]interface{}{
				"embedding": []float64{0.1, 0.2, 0.3, 0.4, 0.5},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case "/api/tags":
			// Health check endpoint
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"models":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Set environment variable for mock server
	originalHost := os.Getenv("OLLAMA_HOST")
	os.Setenv("OLLAMA_HOST", server.URL)
	defer func() {
		if originalHost != "" {
			os.Setenv("OLLAMA_HOST", originalHost)
		} else {
			os.Unsetenv("OLLAMA_HOST")
		}
	}()

	golden.Run(t, "tests/golden", ".txt", ".jsonl", func(srcPath string) ([]byte, error) {
		return runWaferOnFile(srcPath)
	})
}

func runWaferOnFile(srcPath string) ([]byte, error) {
	// Create temporary output file
	tmpDir, err := os.MkdirTemp("", "wafer-golden-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "output.jsonl")
	inputDir := filepath.Dir(srcPath)

	// Build wafer binary if it doesn't exist
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Determine binary name based on platform
	binaryName := "wafer"
	if runtime.GOOS == "windows" {
		binaryName = "wafer.exe"
	}

	waferBin := filepath.Join(wd, binaryName)
	if _, err := os.Stat(waferBin); os.IsNotExist(err) {
		// Build wafer using go build from the repo root
		// Find the repo root by going up from the current working directory
		repoRoot := wd
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
				break
			}
			parent := filepath.Dir(repoRoot)
			if parent == repoRoot {
				return nil, fmt.Errorf("go.mod not found")
			}
			repoRoot = parent
		}

		cmd := exec.Command("go", "build", "-o", waferBin, "./cmd/wafer")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to build wafer: %w", err)
		}
	}

	// Run wafer ingest command
	cmd := exec.Command(waferBin, "ingest", inputDir,
		"--output", outputPath,
		"--chunk-size", "50", // Small chunk size for testing
		"--model", "test-model")

	// Set environment for the command
	ollamaHost := os.Getenv("OLLAMA_HOST")

	// Create a new environment with OLLAMA_HOST explicitly set
	env := []string{}
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "OLLAMA_HOST=") {
			env = append(env, e)
		}
	}
	env = append(env, "OLLAMA_HOST="+ollamaHost)
	cmd.Env = env

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("wafer command failed: %w\nOutput: %s", err, output)
	}

	// Read the generated JSONL file
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("output file was not created: %s", outputPath)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	// Validate JSONL format
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, fmt.Errorf("invalid JSON on line %d: %w", i+1, err)
		}

		// Validate required fields
		requiredFields := []string{"id", "source_file", "chunk_index", "text", "embedding", "word_count", "created_at"}
		for _, field := range requiredFields {
			if _, exists := record[field]; !exists {
				return nil, fmt.Errorf("missing required field '%s' on line %d", field, i+1)
			}
		}
	}

	return content, nil
}
