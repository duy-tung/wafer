package config

import (
	"path/filepath"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Directory: tmpDir,
				Model:     "test-model",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: 300,
			},
			wantErr: false,
		},
		{
			name: "non-existent directory",
			config: &Config{
				Directory: "/non/existent/path",
				Model:     "test-model",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid chunk size",
			config: &Config{
				Directory: tmpDir,
				Model:     "test-model",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: -1,
			},
			wantErr: true,
		},
		{
			name: "empty model",
			config: &Config{
				Directory: tmpDir,
				Model:     "",
				Output:    filepath.Join(tmpDir, "output.jsonl"),
				ChunkSize: 300,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetAbsolutePath(t *testing.T) {
	config := &Config{
		Directory: ".",
	}
	
	path, err := config.GetAbsolutePath()
	if err != nil {
		t.Errorf("GetAbsolutePath() error = %v", err)
	}
	
	if !filepath.IsAbs(path) {
		t.Errorf("GetAbsolutePath() returned relative path: %s", path)
	}
}

func TestConfig_GetOutputPath(t *testing.T) {
	config := &Config{
		Output: "output.jsonl",
	}
	
	path, err := config.GetOutputPath()
	if err != nil {
		t.Errorf("GetOutputPath() error = %v", err)
	}
	
	if !filepath.IsAbs(path) {
		t.Errorf("GetOutputPath() returned relative path: %s", path)
	}
}
