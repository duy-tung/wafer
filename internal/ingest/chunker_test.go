package ingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChunker_ChunkText(t *testing.T) {
	chunker := NewChunker(10) // Small chunk size for testing

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

func TestChunker_ChunkFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "This is a test file with some content for chunking."

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	chunker := NewChunker(300)
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

func TestChunker_ChunkFile_NonExistent(t *testing.T) {
	chunker := NewChunker(300)
	_, err := chunker.ChunkFile("/non/existent/file.txt")
	if err == nil {
		t.Error("ChunkFile() expected error for non-existent file")
	}
}

func TestNewChunker(t *testing.T) {
	chunkSize := 500
	chunker := NewChunker(chunkSize)
	
	if chunker == nil {
		t.Error("NewChunker() returned nil")
	}
	
	// Test that the chunker works with the specified chunk size
	text := strings.Repeat("word ", chunkSize*2) // Text with 2x chunk size
	chunks := chunker.ChunkText(text)
	
	if len(chunks) < 2 {
		t.Errorf("Expected at least 2 chunks for large text, got %d", len(chunks))
	}
}
