package generator

import (
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name        string
		config      GeneratorConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "valid openai config",
			config: GeneratorConfig{
				Provider:  "openai",
				Model:     "gpt-4o-mini",
				APIKey:    "test-key",
				MaxTokens: 2000,
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: GeneratorConfig{
				Provider:  "openai",
				Model:     "gpt-4o-mini",
				MaxTokens: 2000,
			},
			wantErr:     true,
			errContains: "API key is required",
		},
		{
			name: "unsupported provider",
			config: GeneratorConfig{
				Provider:  "unsupported",
				Model:     "test-model",
				APIKey:    "test-key",
				MaxTokens: 2000,
			},
			wantErr:     true,
			errContains: "unsupported provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenerator(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("NewGenerator() error = nil, wantErr true")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewGenerator() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("NewGenerator() error = %v, wantErr false", err)
				return
			}
			if got == nil {
				t.Error("NewGenerator() returned nil generator")
			}
		})
	}
}
