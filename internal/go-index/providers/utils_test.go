package providers_test

import (
	"testing"

	"github.com/mauhlik/go-index/internal/go-index/providers"
)

func TestExtractVersionFromFilename(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.filename, func(t *testing.T) {
			t.Parallel()

			got := providers.ExtractVersionFromFilename(testCase.filename, testCase.artifactName)
			if got != testCase.expected {
				t.Errorf("ExtractVersionFromFilename(%q, %q) = %v; want %v",
					testCase.filename, testCase.artifactName, got, testCase.expected)
			}
		})
	}
}

func TestContainsNumbers(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()

			got := providers.ContainsNumbers(testCase.input)

			if got != testCase.expected {
				t.Errorf("ContainsNumbers(%q) = %v; want %v", testCase.input, got, testCase.expected)
			}
		})
	}
}
