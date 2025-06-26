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

func TestChunker(t *testing.T) {
	chunker := ingest.NewChunker(10) // Small chunk size for testing

	tests := []struct {
		name      string
		text      string
		wantCount int
	}{
		{
			name:      "empty text",
			text:      "",
			wantCount: 0,
		},
		{
			name:      "single word",
			text:      "hello",
			wantCount: 1,
		},
		{
			name:      "small text under chunk size",
			text:      "hello world this is a test",
			wantCount: 1,
		},
		{
			name:      "text requiring multiple chunks",
			text:      strings.Repeat("word ", 25), // 25 words, should create 3 chunks
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := chunker.ChunkText(tt.text)
			if len(chunks) != tt.wantCount {
				t.Errorf("ChunkText() got %d chunks, want %d", len(chunks), tt.wantCount)
			}

			// Verify chunk indices are sequential
			for i, chunk := range chunks {
				if chunk.Index != i {
					t.Errorf("Chunk %d has index %d, want %d", i, chunk.Index, i)
				}
			}

			// Verify word counts are reasonable
			for _, chunk := range chunks {
				if chunk.WordCount <= 0 {
					t.Errorf("Chunk has invalid word count: %d", chunk.WordCount)
				}
			}
		})
	}
}

func TestChunkerFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "This is a test file with some content for chunking."

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	chunker := ingest.NewChunker(300)
	chunks, err := chunker.ChunkFile(testFile)
	if err != nil {
		t.Fatalf("ChunkFile() error = %v", err)
	}

	if len(chunks) != 1 {
		t.Errorf("ChunkFile() got %d chunks, want 1", len(chunks))
	}

	if chunks[0].Text != content {
		t.Errorf("ChunkFile() got text %q, want %q", chunks[0].Text, content)
	}
}

func TestEmbedder(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			http.NotFound(w, r)
			return
		}

		// Return mock embedding
		response := ingest.EmbeddingResponse{
			Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	embedder := ingest.NewEmbedder("test-model")
	embedder.SetBaseURL(server.URL)

	ctx := context.Background()
	embedding, err := embedder.GetEmbedding(ctx, "test text")
	if err != nil {
		t.Fatalf("GetEmbedding() error = %v", err)
	}

	expectedLen := 5
	if len(embedding) != expectedLen {
		t.Errorf("GetEmbedding() got embedding length %d, want %d", len(embedding), expectedLen)
	}

	expected := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	for i, val := range embedding {
		if val != expected[i] {
			t.Errorf("GetEmbedding() got embedding[%d] = %f, want %f", i, val, expected[i])
		}
	}
}

func TestWriter(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.jsonl")

	writer, err := ingest.NewWriter(outputPath)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}
	defer writer.Close()

	// Test writing a record
	chunk := ingest.Chunk{
		Text:      "test text",
		WordCount: 2,
		Index:     0,
	}
	embedding := []float64{0.1, 0.2, 0.3}

	err = writer.WriteRecord("test.txt", chunk, embedding)
	if err != nil {
		t.Fatalf("WriteRecord() error = %v", err)
	}

	// Close writer to flush data
	writer.Close()

	// Read and verify the output
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 1 {
		t.Errorf("Expected 1 line in output, got %d", len(lines))
	}

	// Parse the JSON
	var record ingest.VectorRecord
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify fields
	if record.SourceFile != "test.txt" {
		t.Errorf("Expected source_file 'test.txt', got '%s'", record.SourceFile)
	}
	if record.ChunkIndex != 0 {
		t.Errorf("Expected chunk_index 0, got %d", record.ChunkIndex)
	}
	if record.Text != "test text" {
		t.Errorf("Expected text 'test text', got '%s'", record.Text)
	}
	if record.WordCount != 2 {
		t.Errorf("Expected word_count 2, got %d", record.WordCount)
	}
	if len(record.Embedding) != 3 {
		t.Errorf("Expected embedding length 3, got %d", len(record.Embedding))
	}
	if record.ID == "" {
		t.Error("Expected non-empty ID")
	}
	if record.CreatedAt == "" {
		t.Error("Expected non-empty CreatedAt")
	}
}

func TestConfig(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Directory: tmpDir,
				Model:     "test-model",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: 300,
			},
			wantErr: false,
		},
		{
			name: "non-existent directory",
			config: &config.Config{
				Directory: "/non/existent/path",
				Model:     "test-model",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid chunk size",
			config: &config.Config{
				Directory: tmpDir,
				Model:     "test-model",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: -1,
			},
			wantErr: true,
		},
		{
			name: "empty model",
			config: &config.Config{
				Directory: tmpDir,
				Model:     "",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: 300,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
