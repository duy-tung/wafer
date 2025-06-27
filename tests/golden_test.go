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
	"sync"
	"testing"

	"wafer/golden"
)

var buildWaferOnce sync.Once

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
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case "/api/tags":
			// Health check endpoint
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"models":[]}`)); err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Set environment variable for mock server
	t.Setenv("OLLAMA_HOST", server.URL)

	golden.Run(t, "tests/golden", ".txt", ".jsonl", func(srcPath string) ([]byte, error) {
		return runWaferOnFile(t, srcPath)
	})
}

func runWaferOnFile(t *testing.T, srcPath string) ([]byte, error) {
	// Create temporary output file
	tmpDir, err := os.MkdirTemp("", "wafer-golden-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "output.jsonl")

	// Create a temporary directory with just the single test file
	testInputDir := filepath.Join(tmpDir, "input")
	if err := os.MkdirAll(testInputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create test input dir: %w", err)
	}

	// Copy the source file to the test input directory
	srcContent, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	testFilePath := filepath.Join(testInputDir, filepath.Base(srcPath))
	if err := os.WriteFile(testFilePath, srcContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write test file: %w", err)
	}

	// Build wafer binary if it doesn't exist
	var buildErr error
	buildWaferOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			buildErr = fmt.Errorf("failed to get working directory: %w", err)
			return
		}

		// Determine binary name based on platform
		binaryName := "wafer"
		if runtime.GOOS == "windows" {
			binaryName = "wafer.exe"
		}

		// Find the repo root by going up from the current working directory
		repoRoot := wd
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
				break
			}
			parent := filepath.Dir(repoRoot)
			if parent == repoRoot {
				buildErr = fmt.Errorf("go.mod not found")
				return
			}
			repoRoot = parent
		}

		waferBin := filepath.Join(repoRoot, "bin", binaryName)
		if err := os.MkdirAll(filepath.Dir(waferBin), 0755); err != nil {
			buildErr = fmt.Errorf("failed to create bin directory: %w", err)
			return
		}

		buildCmd := exec.Command("go", "build", "-o", waferBin, "./cmd/wafer")
		buildCmd.Dir = repoRoot
		if output, err := buildCmd.CombinedOutput(); err != nil {
			buildErr = fmt.Errorf("failed to build wafer: %w\nBuild output: %s", err, output)
			return
		}

		// Ensure the binary has execute permissions on Unix systems
		if runtime.GOOS != "windows" {
			if err := os.Chmod(waferBin, 0755); err != nil {
				buildErr = fmt.Errorf("failed to set execute permissions: %w", err)
				return
			}
		}
	})

	if buildErr != nil {
		return nil, buildErr
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
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

	binaryName := "wafer"
	if runtime.GOOS == "windows" {
		binaryName = "wafer.exe"
	}
	waferBin := filepath.Join(repoRoot, "bin", binaryName)

	// Run wafer ingest command on the test input directory
	cmd := exec.Command(waferBin, "ingest", testInputDir,
		"--output", outputPath,
		"--chunk-size", "50", // Small chunk size for testing
		"--model", "test-model")

	// Set environment for the command
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		return nil, fmt.Errorf("OLLAMA_HOST environment variable not set")
	}

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
