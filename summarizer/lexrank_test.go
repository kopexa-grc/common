package summarizer

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewLexRankSummarizer(t *testing.T) {
	tests := []struct {
		name          string
		maxSentences  int
		expectError   bool
		expectedError error
	}{
		{
			name:         "valid maxSentences",
			maxSentences: 5,
			expectError:  false,
		},
		{
			name:         "maxSentences is 1",
			maxSentences: 1,
			expectError:  false,
		},
		{
			name:          "maxSentences is 0",
			maxSentences:  0,
			expectError:   true,
			expectedError: ErrInvalidMaxSentences,
		},
		{
			name:          "maxSentences is negative",
			maxSentences:  -1,
			expectError:   true,
			expectedError: ErrInvalidMaxSentences,
		},
		{
			name:          "maxSentences is very negative",
			maxSentences:  -100,
			expectError:   true,
			expectedError: ErrInvalidMaxSentences,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summarizer, err := newLexRankSummarizer(tt.maxSentences)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}

				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}

				if summarizer != nil {
					t.Errorf("Expected nil summarizer when error occurs, got %v", summarizer)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if summarizer == nil {
					t.Errorf("Expected summarizer but got nil")
					return
				}

				if summarizer.maxSentences != tt.maxSentences {
					t.Errorf("Expected maxSentences %d, got %d", tt.maxSentences, summarizer.maxSentences)
				}
			}
		})
	}
}

func TestLexRankSummarizer_Summarize(t *testing.T) {
	// Create a valid summarizer for testing
	summarizer, err := newLexRankSummarizer(3)
	if err != nil {
		t.Fatalf("Failed to create summarizer: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expectError bool
		checkResult func(string) bool
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   \t\n   ",
			expectError: true,
		},
		{
			name:        "single sentence",
			input:       "This is a single sentence.",
			expectError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "This is a single sentence")
			},
		},
		{
			name: "multiple sentences",
			input: `This is the first sentence. This is the second sentence. 
					This is the third sentence. This is the fourth sentence. 
					This is the fifth sentence.`,
			expectError: false,
			checkResult: func(result string) bool {
				// Should return a summary with fewer sentences than input
				sentences := strings.Split(result, ". ")
				return len(sentences) <= 3 && len(result) > 0
			},
		},
		{
			name:        "very short text",
			input:       "Hi.",
			expectError: false,
			checkResult: func(result string) bool {
				return result == "Hi."
			},
		},
		{
			name:        "text with special characters",
			input:       "Hello! How are you? I'm doing well. This is great!",
			expectError: false,
			checkResult: func(result string) bool {
				return len(result) > 0 && (strings.Contains(result, "Hello") || strings.Contains(result, "How are you") || strings.Contains(result, "doing well") || strings.Contains(result, "great"))
			},
		},
		{
			name:        "deutscher Text",
			input:       "Hallo! Wie geht es dir? Mir geht es gut. Das ist großartig!",
			expectError: false,
			checkResult: func(result string) bool {
				return len(result) > 0 && (strings.Contains(result, "Hallo") || strings.Contains(result, "Wie geht es dir") || strings.Contains(result, "Mir geht es gut") || strings.Contains(result, "großartig"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := summarizer.Summarize(ctx, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}

				if !errors.Is(err, ErrSentenceEmpty) {
					t.Errorf("Expected ErrSentenceEmpty, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if tt.checkResult != nil && !tt.checkResult(result) {
					t.Errorf("Result validation failed for input: %q, result: %q", tt.input, result)
				}
			}
		})
	}
}

func TestLexRankSummarizer_SummarizeWithContext(t *testing.T) {
	summarizer, err := newLexRankSummarizer(2)
	if err != nil {
		t.Fatalf("Failed to create summarizer: %v", err)
	}

	// Test context cancellation
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := summarizer.Summarize(ctx, "This is a test sentence.")
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})

	// Test context timeout
	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait a bit to ensure timeout
		time.Sleep(1 * time.Millisecond)

		_, err := summarizer.Summarize(ctx, "This is a test sentence.")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected context.DeadlineExceeded, got %v", err)
		}
	})

	// Test normal context
	t.Run("normal context", func(t *testing.T) {
		ctx := context.Background()

		result, err := summarizer.Summarize(ctx, "This is a test sentence.")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}
	})
}

func TestLexRankSummarizer_GetMaxSentences(t *testing.T) {
	tests := []struct {
		name         string
		maxSentences int
	}{
		{"one sentence", 1},
		{"three sentences", 3},
		{"ten sentences", 10},
		{"large number", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summarizer, err := newLexRankSummarizer(tt.maxSentences)
			if err != nil {
				t.Fatalf("Failed to create summarizer: %v", err)
			}

			result := summarizer.GetMaxSentences()
			if result != tt.maxSentences {
				t.Errorf("Expected %d, got %d", tt.maxSentences, result)
			}
		})
	}
}

func TestLexRankSummarizer_EdgeCases(t *testing.T) {
	summarizer, err := newLexRankSummarizer(1)
	if err != nil {
		t.Fatalf("Failed to create summarizer: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "single word",
			input: "Hello",
		},
		{
			name:  "single word with period",
			input: "Hello.",
		},
		{
			name:  "multiple spaces",
			input: "Hello    world.   How   are   you?",
		},
		{
			name:  "text with newlines",
			input: "Hello\nworld.\nHow\nare\nyou?",
		},
		{
			name:  "text with tabs",
			input: "Hello\tworld.\tHow\tare\tyou?",
		},
		{
			name:  "text with mixed whitespace",
			input: "Hello \t\n world. \n\t How \t are \n you?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := summarizer.Summarize(ctx, tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("Expected non-empty result")
			}

			// Result should not be longer than input (after trimming)
			trimmedInput := strings.TrimSpace(tt.input)
			if len(result) > len(trimmedInput) {
				t.Errorf("Result is longer than input: result=%d, input=%d", len(result), len(trimmedInput))
			}
		})
	}
}

func TestLexRankSummarizer_Performance(t *testing.T) {
	// Test with a larger text to ensure performance is reasonable
	summarizer, err := newLexRankSummarizer(5)
	if err != nil {
		t.Fatalf("Failed to create summarizer: %v", err)
	}

	// Create a longer text for testing
	longText := strings.Repeat("This is a test sentence. ", 50)

	ctx := context.Background()
	start := time.Now()

	result, err := summarizer.Summarize(ctx, longText)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Performance check: should complete within reasonable time (1 second)
	if duration > time.Second {
		t.Errorf("Summarization took too long: %v", duration)
	}

	// Result should be significantly shorter than input
	if len(result) >= len(longText) {
		t.Errorf("Summary should be shorter than input: summary=%d, input=%d", len(result), len(longText))
	}
}

func BenchmarkLexRankSummarizer_Summarize(b *testing.B) {
	summarizer, err := newLexRankSummarizer(3)
	if err != nil {
		b.Fatalf("Failed to create summarizer: %v", err)
	}

	testText := `This is the first sentence of the benchmark test. 
	This is the second sentence that provides more context. 
	This is the third sentence that adds additional information. 
	This is the fourth sentence that continues the narrative. 
	This is the fifth sentence that concludes the test text.`

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := summarizer.Summarize(ctx, testText)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}
