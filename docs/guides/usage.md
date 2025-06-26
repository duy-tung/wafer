# Usage Guide

This guide provides comprehensive information on how to use the wafer CLI tool effectively.

## Basic Usage

### Command Structure

```bash
wafer ingest <directory_path> [OPTIONS]
```

### Quick Start

Process all `.txt` files in a directory:

```bash
wafer ingest ./documents
```

This will:
1. Recursively find all `.txt` files in `./documents`
2. Split each file into ~300 word chunks
3. Generate embeddings using the `nomic-embed-text` model
4. Save results to `storage/vectors.jsonl`

## Command Options

### Required Arguments

- `<directory_path>`: Path to the directory containing `.txt` files to process

### Optional Flags

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `--model` | Ollama model name | `nomic-embed-text` | `--model=all-minilm` |
| `--output` | Output file path | `storage/vectors.jsonl` | `--output=./embeddings.jsonl` |
| `--chunk-size` | Words per chunk | `300` | `--chunk-size=500` |

### Global Flags

| Flag | Description |
|------|-------------|
| `--version` | Show version information |
| `--help` | Show help message |

## Examples

### Basic Examples

**Process a single directory:**
```bash
wafer ingest ./my-documents
```

**Use a different model:**
```bash
wafer ingest ./documents --model=mxbai-embed-large
```

**Custom output location:**
```bash
wafer ingest ./documents --output=./my-embeddings.jsonl
```

**Larger chunks for better context:**
```bash
wafer ingest ./documents --chunk-size=500
```

### Advanced Examples

**Research paper processing:**
```bash
wafer ingest ./research-papers \
  --model=mxbai-embed-large \
  --output=./research-embeddings.jsonl \
  --chunk-size=400
```

**Quick processing with smaller chunks:**
```bash
wafer ingest ./notes \
  --model=all-minilm \
  --chunk-size=200 \
  --output=./quick-embeddings.jsonl
```

**Processing with custom Ollama host:**
```bash
OLLAMA_HOST=http://remote-server:11434 \
wafer ingest ./documents
```

### Docker Examples

**Basic Docker usage:**
```bash
# Process documents with Docker
docker run --rm \
  -v $(pwd)/documents:/data \
  -v $(pwd)/output:/app/storage \
  --network host \
  ghcr.io/duy-tung/wafer:latest \
  ingest /data
```

**Docker with custom settings:**
```bash
# Research papers with Docker
docker run --rm \
  -v $(pwd)/research-papers:/data \
  -v $(pwd)/embeddings:/app/storage \
  --network host \
  ghcr.io/duy-tung/wafer:latest \
  ingest /data --model=mxbai-embed-large --chunk-size=400 --output=/app/storage/research.jsonl
```

**Docker Compose workflow:**
```yaml
version: '3.8'
services:
  wafer:
    image: ghcr.io/duy-tung/wafer:latest
    volumes:
      - ./documents:/data:ro
      - ./output:/app/storage
    command: ingest /data --chunk-size=300
    network_mode: host
    restart: "no"
```

## Understanding the Output

### Output Format

Wafer generates JSONL (JSON Lines) format where each line is a valid JSON object:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "source_file": "documents/example.txt",
  "chunk_index": 0,
  "text": "This is the actual text content...",
  "embedding": [0.1234, -0.5678, 0.9012, ...],
  "word_count": 299,
  "created_at": "2024-01-15T10:30:45Z"
}
```

### Field Descriptions

- **id**: Unique identifier for each chunk (UUID v4)
- **source_file**: Relative path from the input directory
- **chunk_index**: Sequential number of the chunk within the file (0-based)
- **text**: The actual text content of the chunk
- **embedding**: Array of floating-point numbers representing the embedding
- **word_count**: Actual number of words in this chunk
- **created_at**: ISO 8601 timestamp when the record was created

### Reading the Output

**Using jq to explore the output:**
```bash
# Count total records
cat storage/vectors.jsonl | wc -l

# View first record
head -n1 storage/vectors.jsonl | jq .

# List all source files
cat storage/vectors.jsonl | jq -r '.source_file' | sort | uniq

# Find chunks from a specific file
cat storage/vectors.jsonl | jq 'select(.source_file == "documents/example.txt")'

# Get embedding dimensions
head -n1 storage/vectors.jsonl | jq '.embedding | length'
```

**Using Python to process the output:**
```python
import json

# Read all records
records = []
with open('storage/vectors.jsonl', 'r') as f:
    for line in f:
        records.append(json.loads(line))

print(f"Total records: {len(records)}")
print(f"Embedding dimensions: {len(records[0]['embedding'])}")

# Group by source file
from collections import defaultdict
by_file = defaultdict(list)
for record in records:
    by_file[record['source_file']].append(record)

for file, chunks in by_file.items():
    print(f"{file}: {len(chunks)} chunks")
```

## File Processing Behavior

### File Discovery

- Recursively searches all subdirectories
- Only processes files with `.txt` extension (case-insensitive)
- Skips files that cannot be read (logs warnings)
- Processes files in alphabetical order

### Text Processing

- Reads files as UTF-8 encoded text
- Splits text into chunks at word boundaries
- Preserves original text formatting within chunks
- Handles Unicode characters properly
- Filters out tokens that don't contain letters or digits

### Chunking Logic

- Target chunk size is approximate (±10% variation is normal)
- Preserves word boundaries (never splits words)
- Files smaller than chunk size become single chunks
- Empty files or files with no valid words are skipped
- Chunk indices are sequential within each file

## Error Handling

### Common Scenarios

**File Access Issues:**
- Permission denied → File skipped with warning
- File not found → Directory traversal continues
- Corrupted files → File skipped with error log

**API Issues:**
- Ollama not running → Process stops with error
- Network timeout → Retries with exponential backoff
- Invalid model → Process stops with error

**Output Issues:**
- Output directory doesn't exist → Created automatically
- Insufficient disk space → Process stops with error
- Permission denied on output → Process stops with error

### Exit Codes

- `0`: Success (all files processed)
- `1`: Critical error (configuration, API unavailable)
- `2`: Partial success (some files failed)

## Performance Considerations

### Processing Speed

Factors affecting speed:
- **Model size**: Larger models are slower but more accurate
- **Chunk size**: Larger chunks take longer to process
- **Network latency**: Local Ollama is fastest
- **File size**: Larger files create more chunks

Typical performance:
- Small model (all-minilm): ~500 chunks/minute
- Medium model (nomic-embed-text): ~200 chunks/minute
- Large model (mxbai-embed-large): ~100 chunks/minute

### Memory Usage

- Base application: ~10MB
- Per chunk: ~1KB temporary memory
- Output buffering: Minimal (writes immediately)
- Ollama model: Varies by model (500MB - 4GB)

### Disk Usage

- Output file size: ~2-5x input text size
- Temporary files: None created
- Log files: Not created by default

## Best Practices

### Choosing Chunk Size

- **Small chunks (100-200 words)**: Better for precise search, more records
- **Medium chunks (300-400 words)**: Good balance, default recommendation
- **Large chunks (500+ words)**: Better context, fewer records

### Model Selection

- **nomic-embed-text**: Best general-purpose choice (default)
- **all-minilm**: Fastest, good for large datasets
- **mxbai-embed-large**: Highest quality, slower processing

### Directory Organization

```
project/
├── documents/           # Input directory
│   ├── research/
│   ├── notes/
│   └── references/
├── embeddings/         # Output directory
│   ├── research.jsonl
│   ├── notes.jsonl
│   └── references.jsonl
└── scripts/           # Processing scripts
```

### Batch Processing

Process different document types separately:

```bash
# Process research papers with large chunks
wafer ingest ./documents/research \
  --chunk-size=500 \
  --output=./embeddings/research.jsonl

# Process notes with smaller chunks
wafer ingest ./documents/notes \
  --chunk-size=200 \
  --output=./embeddings/notes.jsonl
```

## Integration Examples

### Vector Database Integration

**Pinecone:**
```python
import json
import pinecone

# Initialize Pinecone
pinecone.init(api_key="your-api-key", environment="your-env")
index = pinecone.Index("your-index")

# Upload vectors
with open('storage/vectors.jsonl', 'r') as f:
    for line in f:
        record = json.loads(line)
        index.upsert([(
            record['id'],
            record['embedding'],
            {
                'source_file': record['source_file'],
                'text': record['text'],
                'chunk_index': record['chunk_index']
            }
        )])
```

**Weaviate:**
```python
import json
import weaviate

client = weaviate.Client("http://localhost:8080")

with open('storage/vectors.jsonl', 'r') as f:
    for line in f:
        record = json.loads(line)
        client.data_object.create(
            data_object={
                'text': record['text'],
                'source_file': record['source_file'],
                'chunk_index': record['chunk_index']
            },
            class_name="Document",
            vector=record['embedding']
        )
```

## Troubleshooting

See the [Installation Guide](installation.md#troubleshooting) for common issues and solutions.
