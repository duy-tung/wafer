package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alexflint/go-arg"

	"wafer/internal/config"
	"wafer/internal/ingest"
)

var (
	version   = "dev"
	gitCommit = "unknown"
	buildTime = "unknown"
)

type CLI struct {
	Ingest  *IngestCmd `arg:"subcommand:ingest" help:"Process text files and generate embeddings"`
	Version bool       `arg:"--version" help:"Print version info and exit"`
}

type IngestCmd struct {
	Directory string `arg:"positional,required" help:"Directory path to process"`
	Model     string `arg:"--model" help:"Ollama model name" default:"nomic-embed-text"`
	Output    string `arg:"--output" help:"Output file path" default:"storage/vectors.jsonl"`
	ChunkSize int    `arg:"--chunk-size" help:"Chunk size in words" default:"300"`
}

func main() {
	// Set up structured logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	var cli CLI
	parser := arg.MustParse(&cli)

	if cli.Version {
		printVersion()
		return
	}

	if cli.Ingest == nil {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	// Create configuration
	cfg := &config.Config{
		Directory: cli.Ingest.Directory,
		Model:     cli.Ingest.Model,
		Output:    cli.Ingest.Output,
		ChunkSize: cli.Ingest.ChunkSize,
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		slog.Error("Configuration error", "error", err)
		os.Exit(1)
	}

	// Run the ingest process
	processor := ingest.NewProcessor(cfg)
	if err := processor.Process(); err != nil {
		slog.Error("Processing failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Processing completed successfully")
}

func printVersion() {
	fmt.Printf("wafer v%s (%s, built %s)\n", version, gitCommit, buildTime)
}
