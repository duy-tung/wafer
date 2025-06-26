# Embedding Features

This document provides detailed information about wafer's embedding generation capabilities and Ollama integration.

## Overview

Wafer uses [Ollama](https://ollama.ai/) to generate high-quality text embeddings. Embeddings are dense vector representations of text that capture semantic meaning, making them useful for search, similarity comparison, and machine learning applications.

## Supported Models

### Default Model: nomic-embed-text

**Characteristics:**
- **Dimensions**: 768
- **Context Length**: 8192 tokens
- **Performance**: Excellent balance of speed and quality
- **Use Cases**: General-purpose text embedding

**Installation:**
```bash
ollama pull nomic-embed-text
```

### Alternative Models

#### all-minilm
- **Dimensions**: 384
- **Context Length**: 512 tokens
- **Performance**: Fast, lightweight
- **Use Cases**: Large datasets, real-time applications

```bash
ollama pull all-minilm
wafer ingest ./documents --model=all-minilm
```

#### mxbai-embed-large
- **Dimensions**: 1024
- **Context Length**: 512 tokens
- **Performance**: High accuracy, slower
- **Use Cases**: High-precision applications, research

```bash
ollama pull mxbai-embed-large
wafer ingest ./documents --model=mxbai-embed-large
```

## Embedding Process

### Text Preprocessing

1. **UTF-8 Decoding**: Files are read as UTF-8 text
2. **Chunking**: Text is split into configurable word-count chunks
3. **Word Boundary Preservation**: Chunks never split words
4. **Whitespace Normalization**: Extra whitespace is cleaned up

### API Integration

**Request Format:**
```json
{
  "model": "nomic-embed-text",
  "prompt": "Your text chunk here..."
}
```

**Response Format:**
```json
{
  "embedding": [0.1234, -0.5678, 0.9012, ...]
}
```

### Error Handling

- **Retry Logic**: 3 attempts with exponential backoff
- **Timeout**: 30-second HTTP timeout
- **Rate Limiting**: Automatic backoff on API limits
- **Network Issues**: Graceful handling of connectivity problems

## Configuration Options

### Model Selection

Choose models based on your requirements:

```bash
# Speed-optimized
wafer ingest ./docs --model=all-minilm

# Balanced (default)
wafer ingest ./docs --model=nomic-embed-text

# Quality-optimized
wafer ingest ./docs --model=mxbai-embed-large
```

### Chunk Size Optimization

**Small Chunks (100-200 words):**
- Better for precise search
- More granular results
- Higher processing overhead

```bash
wafer ingest ./docs --chunk-size=150
```

**Medium Chunks (300-400 words):**
- Good balance of context and precision
- Recommended for most use cases
- Default setting

```bash
wafer ingest ./docs --chunk-size=300  # Default
```

**Large Chunks (500+ words):**
- Better semantic context
- Fewer total embeddings
- Better for document-level similarity

```bash
wafer ingest ./docs --chunk-size=600
```

### Custom Ollama Configuration

**Remote Ollama Server:**
```bash
OLLAMA_HOST=http://your-server:11434 wafer ingest ./docs
```

**Custom Port:**
```bash
OLLAMA_HOST=http://localhost:8080 wafer ingest ./docs
```

## Output Analysis

### Embedding Properties

**Dimensionality:**
```bash
# Check embedding dimensions
head -n1 storage/vectors.jsonl | jq '.embedding | length'
```

**Value Range:**
- Most models produce embeddings in the range [-1, 1]
- Values are typically normalized
- Exact range depends on the model

**Similarity Calculation:**
```python
import json
import numpy as np

# Load embeddings
records = []
with open('storage/vectors.jsonl', 'r') as f:
    for line in f:
        records.append(json.loads(line))

# Calculate cosine similarity
def cosine_similarity(a, b):
    return np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))

# Compare first two embeddings
emb1 = np.array(records[0]['embedding'])
emb2 = np.array(records[1]['embedding'])
similarity = cosine_similarity(emb1, emb2)
print(f"Similarity: {similarity:.4f}")
```

### Quality Assessment

**Semantic Coherence:**
```python
# Find most similar chunks
similarities = []
base_embedding = np.array(records[0]['embedding'])

for i, record in enumerate(records[1:], 1):
    emb = np.array(record['embedding'])
    sim = cosine_similarity(base_embedding, emb)
    similarities.append((i, sim, record['text'][:100]))

# Sort by similarity
similarities.sort(key=lambda x: x[1], reverse=True)

print("Most similar chunks:")
for idx, sim, text in similarities[:5]:
    print(f"Similarity: {sim:.4f} - {text}...")
```

## Performance Optimization

### Model Performance Comparison

| Model | Dimensions | Speed | Quality | Memory |
|-------|------------|-------|---------|--------|
| all-minilm | 384 | Fast | Good | Low |
| nomic-embed-text | 768 | Medium | Excellent | Medium |
| mxbai-embed-large | 1024 | Slow | Best | High |

### Processing Speed Tips

1. **Use appropriate model for your needs**
2. **Optimize chunk size for your use case**
3. **Run Ollama locally for best performance**
4. **Use SSD storage for better I/O**
5. **Ensure sufficient RAM for the model**

### Batch Processing Strategies

**Process by document type:**
```bash
# Research papers (need high quality)
wafer ingest ./research --model=mxbai-embed-large --chunk-size=500

# Quick notes (speed matters)
wafer ingest ./notes --model=all-minilm --chunk-size=200

# General documents (balanced)
wafer ingest ./general --model=nomic-embed-text --chunk-size=300
```

## Use Cases

### Semantic Search

```python
import json
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

# Load embeddings
embeddings = []
texts = []
with open('storage/vectors.jsonl', 'r') as f:
    for line in f:
        record = json.loads(line)
        embeddings.append(record['embedding'])
        texts.append(record['text'])

embeddings = np.array(embeddings)

# Search function
def search(query_embedding, top_k=5):
    similarities = cosine_similarity([query_embedding], embeddings)[0]
    top_indices = np.argsort(similarities)[::-1][:top_k]
    
    results = []
    for idx in top_indices:
        results.append({
            'text': texts[idx],
            'similarity': similarities[idx]
        })
    return results
```

### Document Clustering

```python
from sklearn.cluster import KMeans
import matplotlib.pyplot as plt

# Cluster embeddings
kmeans = KMeans(n_clusters=5, random_state=42)
clusters = kmeans.fit_predict(embeddings)

# Analyze clusters
for i in range(5):
    cluster_texts = [texts[j] for j in range(len(texts)) if clusters[j] == i]
    print(f"Cluster {i}: {len(cluster_texts)} documents")
    print(f"Sample: {cluster_texts[0][:100]}...")
    print()
```

### Similarity Analysis

```python
# Find duplicate or near-duplicate content
threshold = 0.95
duplicates = []

for i in range(len(embeddings)):
    for j in range(i+1, len(embeddings)):
        sim = cosine_similarity([embeddings[i]], [embeddings[j]])[0][0]
        if sim > threshold:
            duplicates.append((i, j, sim, texts[i][:50], texts[j][:50]))

print(f"Found {len(duplicates)} potential duplicates")
for i, j, sim, text1, text2 in duplicates[:5]:
    print(f"Similarity: {sim:.4f}")
    print(f"Text 1: {text1}...")
    print(f"Text 2: {text2}...")
    print()
```

## Troubleshooting

### Common Issues

**1. "Model not found" error**
```bash
# Check available models
ollama list

# Pull the required model
ollama pull nomic-embed-text
```

**2. "Connection refused" error**
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama if not running
ollama serve
```

**3. "Embedding generation failed" error**
- Check Ollama logs: `ollama logs`
- Verify model is loaded: `ollama ps`
- Check available memory: `free -h`

**4. Slow processing**
- Use a smaller/faster model
- Reduce chunk size
- Check system resources
- Ensure Ollama has sufficient RAM

### Performance Monitoring

**Check Ollama status:**
```bash
# List running models
ollama ps

# Check model info
ollama show nomic-embed-text

# Monitor resource usage
htop  # or top on macOS
```

**Monitor wafer processing:**
```bash
# Enable debug logging
WAFER_LOG_LEVEL=debug wafer ingest ./docs

# Monitor output file growth
watch -n 1 'wc -l storage/vectors.jsonl'
```

### Advanced Configuration

**Custom timeout settings:**
```bash
# Increase timeout for slow models
OLLAMA_REQUEST_TIMEOUT=60 wafer ingest ./docs
```

**Memory optimization:**
```bash
# Process in smaller batches
find ./docs -name "*.txt" | head -100 | xargs -I {} dirname {} | sort -u | head -10 | while read dir; do
    wafer ingest "$dir" --output="embeddings/$(basename "$dir").jsonl"
done
```

## Integration Examples

### Vector Database Upload

See the [Usage Guide](../guides/usage.md#integration-examples) for detailed integration examples with popular vector databases.

### Custom Processing Pipeline

```python
import json
import subprocess
import tempfile
import os

def process_documents_with_custom_preprocessing(docs_dir, output_file):
    """Custom preprocessing before embedding generation"""
    
    # Create temporary directory for preprocessed files
    with tempfile.TemporaryDirectory() as temp_dir:
        # Custom preprocessing (example: remove headers/footers)
        for root, dirs, files in os.walk(docs_dir):
            for file in files:
                if file.endswith('.txt'):
                    input_path = os.path.join(root, file)
                    output_path = os.path.join(temp_dir, file)
                    
                    # Custom preprocessing logic here
                    with open(input_path, 'r') as f:
                        content = f.read()
                    
                    # Remove common headers/footers
                    processed_content = preprocess_text(content)
                    
                    with open(output_path, 'w') as f:
                        f.write(processed_content)
        
        # Run wafer on preprocessed files
        subprocess.run([
            'wafer', 'ingest', temp_dir,
            '--output', output_file,
            '--chunk-size', '400'
        ])

def preprocess_text(text):
    """Custom text preprocessing"""
    lines = text.split('\n')
    # Remove headers, footers, page numbers, etc.
    # This is just an example - customize for your needs
    processed_lines = [line for line in lines if not line.strip().isdigit()]
    return '\n'.join(processed_lines)
```
