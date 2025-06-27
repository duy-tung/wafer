package ingest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEmbedder_GetEmbedding(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			http.NotFound(w, r)
			return
		}

		// Return mock embedding
		response := EmbeddingResponse{
			Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	embedder := NewEmbedder("test-model")
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

func TestEmbedder_HealthCheck(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"models":[]}`)); err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	embedder := NewEmbedder("test-model")
	embedder.SetBaseURL(server.URL)

	ctx := context.Background()
	err := embedder.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck() error = %v", err)
	}
}

func TestEmbedder_HealthCheck_Failure(t *testing.T) {
	embedder := NewEmbedder("test-model")
	embedder.SetBaseURL("http://localhost:99999") // Non-existent server

	ctx := context.Background()
	err := embedder.HealthCheck(ctx)
	if err == nil {
		t.Error("HealthCheck() expected error for non-existent server")
	}
}

func TestNewEmbedder(t *testing.T) {
	model := "test-model"
	embedder := NewEmbedder(model)

	if embedder == nil {
		t.Error("NewEmbedder() returned nil")
	}
}

func TestNewEmbedder_WithOllamaHost(t *testing.T) {
	// Set OLLAMA_HOST environment variable
	originalHost := os.Getenv("OLLAMA_HOST")
	testHost := "http://test.example.com:8080"
	os.Setenv("OLLAMA_HOST", testHost)
	defer func() {
		if originalHost != "" {
			os.Setenv("OLLAMA_HOST", originalHost)
		} else {
			os.Unsetenv("OLLAMA_HOST")
		}
	}()

	embedder := NewEmbedder("test-model")

	if embedder == nil {
		t.Fatal("NewEmbedder() returned nil")
	}

	// Test that it uses the environment variable
	if embedder.baseURL != testHost {
		t.Errorf("NewEmbedder() baseURL = %s, want %s", embedder.baseURL, testHost)
	}
}

func TestEmbedder_SetBaseURL(t *testing.T) {
	embedder := NewEmbedder("test-model")
	newURL := "http://custom.example.com:8080"

	embedder.SetBaseURL(newURL)

	if embedder.baseURL != newURL {
		t.Errorf("SetBaseURL() baseURL = %s, want %s", embedder.baseURL, newURL)
	}
}
