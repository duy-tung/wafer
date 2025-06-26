package ingest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_WriteRecord(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.jsonl")

	writer, err := NewWriter(outputPath)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}
	defer writer.Close()

	// Test writing a record
	chunk := Chunk{
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
	var record VectorRecord
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

func TestWriter_MultipleRecords(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.jsonl")

	writer, err := NewWriter(outputPath)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}
	defer writer.Close()

	// Write multiple records
	for i := 0; i < 3; i++ {
		chunk := Chunk{
			Text:      "test text",
			WordCount: 2,
			Index:     i,
		}
		embedding := []float64{0.1, 0.2, 0.3}

		err = writer.WriteRecord("test.txt", chunk, embedding)
		if err != nil {
			t.Fatalf("WriteRecord() error = %v", err)
		}
	}

	writer.Close()

	// Read and verify the output
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines in output, got %d", len(lines))
	}
}

func TestNewWriter_InvalidPath(t *testing.T) {
	// Try to create writer with invalid path (on Windows, use an invalid drive)
	invalidPath := "Z:\\invalid\\path\\output.jsonl"
	if os.PathSeparator == '/' {
		// On Unix-like systems, use a path that can't be created
		invalidPath = "/proc/invalid/output.jsonl"
	}

	_, err := NewWriter(invalidPath)
	if err == nil {
		t.Error("NewWriter() expected error for invalid path")
	}
}

func TestWriter_WriteAfterClose(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.jsonl")

	writer, err := NewWriter(outputPath)
	if err != nil {
		t.Fatalf("NewWriter() error = %v", err)
	}

	// Close the writer
	writer.Close()

	// Try to write after close
	chunk := Chunk{
		Text:      "test text",
		WordCount: 2,
		Index:     0,
	}
	embedding := []float64{0.1, 0.2, 0.3}

	err = writer.WriteRecord("test.txt", chunk, embedding)
	if err == nil {
		t.Error("WriteRecord() expected error after close")
	}
}
