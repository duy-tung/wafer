package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// VectorRecord represents a single record in the JSONL output
type VectorRecord struct {
	ID         string    `json:"id"`
	SourceFile string    `json:"source_file"`
	ChunkIndex int       `json:"chunk_index"`
	Text       string    `json:"text"`
	Embedding  []float64 `json:"embedding"`
	WordCount  int       `json:"word_count"`
	CreatedAt  string    `json:"created_at"`
}

// Writer handles writing JSONL output
type Writer struct {
	outputPath string
	file       *os.File
}

// NewWriter creates a new writer for the specified output path
func NewWriter(outputPath string) (*Writer, error) {
	// Ensure the output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open file in append mode
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}

	return &Writer{
		outputPath: outputPath,
		file:       file,
	}, nil
}

// WriteRecord writes a single vector record to the JSONL file
func (w *Writer) WriteRecord(sourceFile string, chunk Chunk, embedding []float64) error {
	// Create the record
	record := VectorRecord{
		ID:         uuid.New().String(),
		SourceFile: sourceFile,
		ChunkIndex: chunk.Index,
		Text:       chunk.Text,
		Embedding:  embedding,
		WordCount:  chunk.WordCount,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	// Write to file with newline
	if _, err := w.file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}
	if _, err := w.file.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// Close closes the writer and flushes any remaining data
func (w *Writer) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// GetOutputPath returns the output file path
func (w *Writer) GetOutputPath() string {
	return w.outputPath
}

// GetFileSize returns the current size of the output file
func (w *Writer) GetFileSize() (int64, error) {
	info, err := w.file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	return info.Size(), nil
}

// Flush ensures all data is written to disk
func (w *Writer) Flush() error {
	return w.file.Sync()
}
