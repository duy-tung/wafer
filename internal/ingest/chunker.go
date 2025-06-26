package ingest

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Chunk represents a text chunk with metadata
type Chunk struct {
	Text      string
	WordCount int
	Index     int
}

// Chunker handles text chunking operations
type Chunker struct {
	chunkSize int
}

// NewChunker creates a new chunker with the specified chunk size
func NewChunker(chunkSize int) *Chunker {
	return &Chunker{
		chunkSize: chunkSize,
	}
}

// ChunkFile reads a file and splits it into chunks
func (c *Chunker) ChunkFile(filePath string) ([]Chunk, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return c.ChunkText(string(content)), nil
}

// ChunkText splits text into chunks of approximately the specified word count
func (c *Chunker) ChunkText(text string) []Chunk {
	// Clean and normalize the text
	text = strings.TrimSpace(text)
	if text == "" {
		return []Chunk{}
	}

	// Split text into words while preserving whitespace information
	words := c.tokenizeWords(text)
	if len(words) == 0 {
		return []Chunk{}
	}

	// If the text has fewer words than chunk size, return as single chunk
	if len(words) <= c.chunkSize {
		return []Chunk{{
			Text:      text,
			WordCount: len(words),
			Index:     0,
		}}
	}

	var chunks []Chunk
	chunkIndex := 0

	for i := 0; i < len(words); i += c.chunkSize {
		end := i + c.chunkSize
		if end > len(words) {
			end = len(words)
		}

		chunkWords := words[i:end]
		chunkText := strings.Join(chunkWords, " ")

		// Clean up extra whitespace
		chunkText = strings.TrimSpace(chunkText)

		if chunkText != "" {
			chunks = append(chunks, Chunk{
				Text:      chunkText,
				WordCount: len(chunkWords),
				Index:     chunkIndex,
			})
			chunkIndex++
		}
	}

	return chunks
}

// tokenizeWords splits text into words while preserving word boundaries
func (c *Chunker) tokenizeWords(text string) []string {
	var words []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		// Only include non-empty words that contain at least one letter or digit
		if c.isValidWord(word) {
			words = append(words, word)
		}
	}

	return words
}

// isValidWord checks if a word is valid (contains at least one alphanumeric character)
func (c *Chunker) isValidWord(word string) bool {
	if word == "" {
		return false
	}

	for _, r := range word {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}

	return false
}
