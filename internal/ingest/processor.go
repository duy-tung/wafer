package ingest

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"wafer/internal/config"
)

// ProcessorStats holds statistics about the processing run
type ProcessorStats struct {
	FilesProcessed int
	FilesSkipped   int
	ChunksCreated  int
	TotalErrors    int
	StartTime      time.Time
	EndTime        time.Time
}

// Processor orchestrates the entire ingestion process
type Processor struct {
	config   *config.Config
	chunker  *Chunker
	embedder *Embedder
	writer   *Writer
	stats    ProcessorStats
}

// NewProcessor creates a new processor with the given configuration
func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		config:   cfg,
		chunker:  NewChunker(cfg.ChunkSize),
		embedder: NewEmbedder(cfg.Model),
		stats:    ProcessorStats{StartTime: time.Now()},
	}
}

// Process runs the complete ingestion workflow
func (p *Processor) Process() error {
	ctx := context.Background()

	slog.Info("Starting wafer ingestion process",
		"directory", p.config.Directory,
		"model", p.config.Model,
		"output", p.config.Output,
		"chunk_size", p.config.ChunkSize)

	// Health check Ollama API
	slog.Info("Checking Ollama API connectivity...")
	if err := p.embedder.HealthCheck(ctx); err != nil {
		return fmt.Errorf("Ollama API health check failed: %w", err)
	}
	slog.Info("Ollama API is accessible")

	// Initialize writer
	writer, err := NewWriter(p.config.Output)
	if err != nil {
		return fmt.Errorf("failed to initialize writer: %w", err)
	}
	defer writer.Close()
	p.writer = writer

	// Discover .txt files
	txtFiles, err := p.discoverTextFiles()
	if err != nil {
		return fmt.Errorf("failed to discover text files: %w", err)
	}

	if len(txtFiles) == 0 {
		slog.Warn("No .txt files found in directory", "directory", p.config.Directory)
		return nil
	}

	slog.Info("Found text files to process", "count", len(txtFiles))

	// Process each file
	for i, filePath := range txtFiles {
		slog.Info("Processing file",
			"file", filePath,
			"progress", fmt.Sprintf("%d/%d", i+1, len(txtFiles)))

		if err := p.processFile(ctx, filePath); err != nil {
			slog.Error("Failed to process file", "file", filePath, "error", err)
			p.stats.FilesSkipped++
			p.stats.TotalErrors++
			continue
		}

		p.stats.FilesProcessed++
	}

	// Finalize
	p.stats.EndTime = time.Now()
	p.printSummary()

	return nil
}

// discoverTextFiles recursively finds all .txt files in the directory
func (p *Processor) discoverTextFiles() ([]string, error) {
	var txtFiles []string

	err := filepath.WalkDir(p.config.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Warn("Error accessing path", "path", path, "error", err)
			return nil // Continue walking
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if it's a .txt file
		if strings.ToLower(filepath.Ext(d.Name())) == ".txt" {
			txtFiles = append(txtFiles, path)
		}

		return nil
	})

	return txtFiles, err
}

// processFile processes a single text file
func (p *Processor) processFile(ctx context.Context, filePath string) error {
	// Get relative path for output
	relPath, err := filepath.Rel(p.config.Directory, filePath)
	if err != nil {
		relPath = filePath // Fallback to absolute path
	}

	// Chunk the file
	chunks, err := p.chunker.ChunkFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to chunk file: %w", err)
	}

	if len(chunks) == 0 {
		slog.Warn("File produced no chunks", "file", filePath)
		return nil
	}

	slog.Debug("File chunked", "file", filePath, "chunks", len(chunks))

	// Process each chunk
	for _, chunk := range chunks {
		if err := p.processChunk(ctx, relPath, chunk); err != nil {
			return fmt.Errorf("failed to process chunk %d: %w", chunk.Index, err)
		}
		p.stats.ChunksCreated++
	}

	return nil
}

// processChunk processes a single text chunk
func (p *Processor) processChunk(ctx context.Context, sourceFile string, chunk Chunk) error {
	// Generate embedding
	embedding, err := p.embedder.GetEmbedding(ctx, chunk.Text)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Write to output
	if err := p.writer.WriteRecord(sourceFile, chunk, embedding); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

// printSummary prints a summary of the processing results
func (p *Processor) printSummary() {
	duration := p.stats.EndTime.Sub(p.stats.StartTime)

	slog.Info("Processing completed",
		"files_processed", p.stats.FilesProcessed,
		"files_skipped", p.stats.FilesSkipped,
		"chunks_created", p.stats.ChunksCreated,
		"total_errors", p.stats.TotalErrors,
		"duration", duration.String(),
		"output_file", p.config.Output)

	if p.stats.TotalErrors > 0 {
		slog.Warn("Processing completed with errors", "error_count", p.stats.TotalErrors)
	}
}
