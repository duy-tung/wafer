package tests

import (
	"strings"
	"testing"

	"wafer/internal/ingest"
)

// BenchmarkChunker tests the performance of text chunking
func BenchmarkChunker(b *testing.B) {
	chunker := ingest.NewChunker(300)

	// Create test text of various sizes
	testCases := []struct {
		name string
		text string
	}{
		{"Small", strings.Repeat("word ", 100)},
		{"Medium", strings.Repeat("word ", 1000)},
		{"Large", strings.Repeat("word ", 10000)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				chunks := chunker.ChunkText(tc.text)
				if len(chunks) == 0 {
					b.Fatal("no chunks produced")
				}
			}
		})
	}
}

// BenchmarkChunkerMemory tests memory allocation during chunking
func BenchmarkChunkerMemory(b *testing.B) {
	chunker := ingest.NewChunker(300)
	text := strings.Repeat("word ", 1000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		chunks := chunker.ChunkText(text)
		_ = chunks // Prevent optimization
	}
}

// BenchmarkChunkerDifferentSizes tests chunking performance with different chunk sizes
func BenchmarkChunkerDifferentSizes(b *testing.B) {
	text := strings.Repeat("word ", 5000)

	chunkSizes := []int{100, 300, 500, 1000}

	for _, size := range chunkSizes {
		b.Run(string(rune(size)), func(b *testing.B) {
			chunker := ingest.NewChunker(size)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				chunks := chunker.ChunkText(text)
				if len(chunks) == 0 {
					b.Fatal("no chunks produced")
				}
			}
		})
	}
}

// BenchmarkJSONLWriting tests the performance of JSONL output writing
func BenchmarkJSONLWriting(b *testing.B) {
	tmpDir := b.TempDir()
	outputPath := tmpDir + "/benchmark.jsonl"

	writer, err := ingest.NewWriter(outputPath)
	if err != nil {
		b.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	chunk := ingest.Chunk{
		Text:      strings.Repeat("benchmark text ", 50),
		WordCount: 100,
		Index:     0,
	}
	embedding := make([]float64, 768) // Typical embedding size
	for i := range embedding {
		embedding[i] = float64(i) / 768.0
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		chunk.Index = i
		err := writer.WriteRecord("benchmark.txt", chunk, embedding)
		if err != nil {
			b.Fatalf("Failed to write record: %v", err)
		}
	}
}

// BenchmarkTextProcessingPipeline tests the complete text processing pipeline
func BenchmarkTextProcessingPipeline(b *testing.B) {
	chunker := ingest.NewChunker(300)
	text := strings.Repeat("This is a benchmark test for the complete text processing pipeline. ", 100)

	tmpDir := b.TempDir()
	outputPath := tmpDir + "/pipeline.jsonl"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Chunk the text
		chunks := chunker.ChunkText(text)

		// Create writer
		writer, err := ingest.NewWriter(outputPath)
		if err != nil {
			b.Fatalf("Failed to create writer: %v", err)
		}

		// Process chunks (simulate embedding with dummy data)
		embedding := make([]float64, 768)
		for j, chunk := range chunks {
			for k := range embedding {
				embedding[k] = float64(j*k) / 1000.0
			}

			err := writer.WriteRecord("benchmark.txt", chunk, embedding)
			if err != nil {
				b.Fatalf("Failed to write record: %v", err)
			}
		}

		writer.Close()
	}
}
