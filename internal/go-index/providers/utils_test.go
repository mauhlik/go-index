package providers

import (
	"testing"
)

func TestExtractVersionFromFilename(t *testing.T) {
	tests := []struct {
		filename     string
		artifactName string
		expected     string
	}{
		{"app-0.0.0.txt", "app", "0.0.0"},
		{"app-0.0.1.txt", "app", "0.0.1"},
		{"app-1.0.0.txt", "app", "1.0.0"},
		{"app-2.0.0.txt", "app", "2.0.0"},
		{"app-2.0.0-beta.txt", "app", "2.0.0-beta"},
		{"app-2.0.0-beta.1.txt", "app", "2.0.0-beta.1"},
		{"app-2.0.0-beta.1", "app", "2.0.0-beta.1"},
		{"app-2.0.0", "app", "2.0.0"},
		{"other-1.0.0.txt", "app", ""},
		{"app-1.0.0", "other", ""},
		{"app", "app", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := ExtractVersionFromFilename(tt.filename, tt.artifactName)
			if got != tt.expected {
				t.Errorf("ExtractVersionFromFilename(%q, %q) = %v; want %v", tt.filename, tt.artifactName, got, tt.expected)
			}
		})
	}
}

func TestContainsNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1.0.0", true},
		{"beta", false},
		{"beta.1", true},
		{"", false},
		{"123", true},
		{"abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := containsNumbers(tt.input)
			if got != tt.expected {
				t.Errorf("containsNumbers(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}
