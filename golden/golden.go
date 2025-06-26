package golden

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update golden files")

// Runner defines a test processor function.
type Runner func(srcPath string) ([]byte, error)

// Run executes golden file-based tests for wafer CLI tool.
func Run(t *testing.T, dir, srcExt, goldenExt string, fn Runner) {
	rootDir := findRepoRoot(t)
	pattern := filepath.Join(rootDir, dir, "*"+srcExt)
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to list %s files in %s: %v", srcExt, dir, err)
	}

	if len(files) == 0 {
		t.Fatalf("no test files found: %s", pattern)
	}

	for _, src := range files {
		name := strings.TrimSuffix(filepath.Base(src), srcExt)
		wantPath := filepath.Join(rootDir, dir, name+goldenExt)

		t.Run(name, func(t *testing.T) {
			start := time.Now()
			got, err := fn(src)
			dur := time.Since(start)

			if err != nil {
				t.Fatalf("process error: %v", err)
			}

			if got == nil {
				t.Fatal("got nil output")
			}

			got = normalizeOutput(rootDir, got)

			if *update {
				if err := os.WriteFile(wantPath, got, 0644); err != nil {
					t.Fatalf("failed to write golden: %v", err)
				}
				t.Logf("updated: %s (processed in %v)", wantPath, dur)
				return
			}

			want, err := os.ReadFile(wantPath)
			if err != nil {
				t.Fatalf("failed to read golden: %v", err)
			}

			want = bytes.TrimSpace(want)
			if !bytes.Equal(got, want) {
				t.Errorf("golden mismatch for %s\n\n--- Got ---\n%s\n\n--- Want ---\n%s\n",
					name+goldenExt, got, want)
				return
			}

			t.Logf("✅ %s (processed in %v)", name, dur)
		})
	}
}

// findRepoRoot walks up to locate the `go.mod` file.
func findRepoRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("cannot determine working directory")
	}

	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatal("go.mod not found (not in Go module)")
	return ""
}

// normalizeOutput strips paths and timing information for clean diffing.
func normalizeOutput(root string, b []byte) []byte {
	out := string(b)

	// Strip absolute paths for stable diffs
	out = strings.ReplaceAll(out, filepath.ToSlash(root)+"/", "")
	out = strings.ReplaceAll(out, filepath.ToSlash(root), "")
	out = strings.ReplaceAll(out, "github.com/duy-tung/wafer/", "")
	out = strings.ReplaceAll(out, "wafer/tests/", "tests/")

	// Remove timing information like "(123ns)" or "(1.0µs)" as durations vary
	// slightly between runs and would cause flaky golden tests
	durRE := regexp.MustCompile(`\([0-9]+(\.[0-9]+)?(ns|µs|ms|s)\)`)
	out = durRE.ReplaceAllString(out, "(X)")

	// Remove timestamps from JSONL output for stable comparisons
	timestampRE := regexp.MustCompile(`"created_at":"[^"]+"`)
	out = timestampRE.ReplaceAllString(out, `"created_at":"2024-01-15T10:30:45Z"`)

	// Remove UUIDs from JSONL output for stable comparisons
	uuidRE := regexp.MustCompile(`"id":"[a-f0-9-]+"`)
	out = uuidRE.ReplaceAllString(out, `"id":"550e8400-e29b-41d4-a716-446655440000"`)

	return []byte(strings.TrimSpace(out))
}
