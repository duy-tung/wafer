package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"wafer/internal/config"
	"wafer/internal/ingest"
)

func TestEndToEndProcessing(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		"file1.txt":        "This is the first test file with some content for processing. It has multiple sentences to test chunking.",
		"file2.txt":        "This is the second test file. It contains different content to verify that multiple files are processed correctly.",
		"subdir/file3.txt": "This file is in a subdirectory to test recursive directory traversal functionality.",
	}

	for relPath, content := range testFiles {
		fullPath := filepath.Join(inputDir, relPath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}

	// Create mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/embeddings":
			// Return mock embedding based on text length
			var req ingest.EmbeddingRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// Generate a simple mock embedding based on text length
			embeddingSize := 5
			embedding := make([]float64, embeddingSize)
			for i := range embedding {
				embedding[i] = float64(len(req.Prompt)%100) / 100.0 // Normalize to 0-1
			}

			response := ingest.EmbeddingResponse{Embedding: embedding}
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

	// Configure the processor
	outputPath := filepath.Join(outputDir, "vectors.jsonl")
	cfg := &config.Config{
		Directory: inputDir,
		Model:     "test-model",
		Output:    outputPath,
		ChunkSize: 50, // Small chunk size for testing
	}

	// Create a custom processor for testing with mock embedder
	embedder := ingest.NewEmbedder(cfg.Model)
	embedder.SetBaseURL(server.URL)

	testProcessor := &TestProcessor{
		config:   cfg,
		chunker:  ingest.NewChunker(cfg.ChunkSize),
		embedder: embedder,
	}

	// Run the processing
	if err := testProcessor.Process(); err != nil {
		t.Fatalf("Processing failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputPath)
	}

	// Read and verify output
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 0 {
		t.Fatal("Output file is empty")
	}

	// Verify each line is valid JSON
	var records []ingest.VectorRecord
	for i, line := range lines {
		if line == "" {
			continue
		}

		var record ingest.VectorRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("Line %d is not valid JSON: %v", i+1, err)
		}

		records = append(records, record)

		// Verify required fields
		if record.ID == "" {
			t.Errorf("Record %d missing ID", i+1)
		}
		if record.SourceFile == "" {
			t.Errorf("Record %d missing SourceFile", i+1)
		}
		if record.Text == "" {
			t.Errorf("Record %d missing Text", i+1)
		}
		if len(record.Embedding) == 0 {
			t.Errorf("Record %d missing Embedding", i+1)
		}
		if record.WordCount <= 0 {
			t.Errorf("Record %d has invalid WordCount: %d", i+1, record.WordCount)
		}
		if record.CreatedAt == "" {
			t.Errorf("Record %d missing CreatedAt", i+1)
		}
	}

	// Verify we processed all files
	sourceFiles := make(map[string]bool)
	for _, record := range records {
		sourceFiles[record.SourceFile] = true
	}

	expectedFiles := []string{"file1.txt", "file2.txt", "subdir/file3.txt"}
	for _, expectedFile := range expectedFiles {
		if !sourceFiles[expectedFile] {
			t.Errorf("Expected file %s not found in output", expectedFile)
		}
	}

	t.Logf("Successfully processed %d records from %d files", len(records), len(sourceFiles))
}

// TestProcessor is a custom processor for testing that allows us to inject a mock embedder
type TestProcessor struct {
	config   *config.Config
	chunker  *ingest.Chunker
	embedder *ingest.Embedder
}

func (p *TestProcessor) Process() error {
	ctx := context.Background()

	// Health check
	if err := p.embedder.HealthCheck(ctx); err != nil {
		return err
	}

	// Initialize writer
	writer, err := ingest.NewWriter(p.config.Output)
	if err != nil {
		return err
	}
	defer writer.Close()

	// Discover files
	var txtFiles []string
	err = filepath.WalkDir(p.config.Directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.ToLower(filepath.Ext(d.Name())) == ".txt" {
			txtFiles = append(txtFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Process each file
	for _, filePath := range txtFiles {
		relPath, _ := filepath.Rel(p.config.Directory, filePath)

		chunks, err := p.chunker.ChunkFile(filePath)
		if err != nil {
			continue
		}

		for _, chunk := range chunks {
			embedding, err := p.embedder.GetEmbedding(ctx, chunk.Text)
			if err != nil {
				continue
			}

			if err := writer.WriteRecord(relPath, chunk, embedding); err != nil {
				continue
			}
		}
	}

	return nil
}
