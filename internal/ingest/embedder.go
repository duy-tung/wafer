package ingest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// EmbeddingRequest represents the request to Ollama API
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse represents the response from Ollama API
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// Embedder handles communication with Ollama API
type Embedder struct {
	client  *http.Client
	baseURL string
	model   string
	retries int
	backoff time.Duration
}

// NewEmbedder creates a new embedder with the specified model
func NewEmbedder(model string) *Embedder {
	return &Embedder{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "http://localhost:11434", // Default Ollama URL
		model:   model,
		retries: 3,
		backoff: time.Second,
	}
}

// SetBaseURL sets the Ollama API base URL
func (e *Embedder) SetBaseURL(url string) {
	e.baseURL = url
}

// GetEmbedding generates an embedding for the given text
func (e *Embedder) GetEmbedding(ctx context.Context, text string) ([]float64, error) {
	var lastErr error

	for attempt := 0; attempt <= e.retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoffDuration := e.backoff * time.Duration(1<<(attempt-1))
			slog.Debug("Retrying embedding request",
				"attempt", attempt,
				"backoff", backoffDuration,
				"text_length", len(text))

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoffDuration):
			}
		}

		embedding, err := e.requestEmbedding(ctx, text)
		if err == nil {
			return embedding, nil
		}

		lastErr = err
		slog.Warn("Embedding request failed",
			"attempt", attempt+1,
			"error", err,
			"text_length", len(text))
	}

	return nil, fmt.Errorf("failed to get embedding after %d attempts: %w", e.retries+1, lastErr)
}

// requestEmbedding makes a single request to the Ollama API
func (e *Embedder) requestEmbedding(ctx context.Context, text string) ([]float64, error) {
	// Prepare request
	reqBody := EmbeddingRequest{
		Model:  e.model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/embeddings", e.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embeddingResp.Embedding) == 0 {
		return nil, fmt.Errorf("received empty embedding")
	}

	return embeddingResp.Embedding, nil
}

// HealthCheck verifies that the Ollama API is accessible
func (e *Embedder) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/tags", e.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama API is not accessible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama API health check failed with status %d", resp.StatusCode)
	}

	return nil
}
