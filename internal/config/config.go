package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the wafer CLI tool
type Config struct {
	Directory string // Directory to process
	Model     string // Ollama model name
	Output    string // Output file path
	ChunkSize int    // Chunk size in words
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Check if directory exists
	if _, err := os.Stat(c.Directory); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", c.Directory)
	}

	// Check if directory is actually a directory
	info, err := os.Stat(c.Directory)
	if err != nil {
		return fmt.Errorf("cannot access directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", c.Directory)
	}

	// Validate chunk size
	if c.ChunkSize <= 0 {
		return fmt.Errorf("chunk size must be positive, got: %d", c.ChunkSize)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(c.Output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}

	// Validate model name is not empty
	if c.Model == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	return nil
}

// GetAbsolutePath returns the absolute path for the directory
func (c *Config) GetAbsolutePath() (string, error) {
	return filepath.Abs(c.Directory)
}

// GetOutputPath returns the absolute path for the output file
func (c *Config) GetOutputPath() (string, error) {
	return filepath.Abs(c.Output)
}
