// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestNewLLMSummarizer(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid OpenAI config",
			config: *NewConfig(
				WithType(TypeLlm),
				WithOpenAI("gpt-4", "test-api-key"),
			),
			expectError: false,
		},
		{
			name: "valid Azure OpenAI config",
			config: *NewConfig(
				WithType(TypeLlm),
				WithOpenAI("gpt-4", "test-api-key",
					WithURL("https://test.openai.azure.com/"),
					WithOption("organization_id", "test-org"),
				),
			),
			expectError: false,
		},
		{
			name: "valid Ollama config",
			config: *NewConfig(
				WithType(TypeLlm),
				WithOllama("llama2", "http://localhost:11434"),
			),
			expectError: false,
		},
		{
			name: "missing LLM config",
			config: *NewConfig(
				WithType(TypeLlm),
			),
			expectError: true,
		},
		{
			name: "LexRank type with LLM config",
			config: *NewConfig(
				WithType(TypeLexrank),
				WithOpenAI("gpt-4", "test-api-key"),
			),
			expectError: false, // Should not error, LLM config is ignored for LexRank
		},
		{
			name: "Azure OpenAI advanced options",
			config: *NewConfig(
				WithType(TypeLlm),
				WithOpenAI("gpt-4", "test-key",
					WithURL("https://julian-local.openai.azure.com/"),
					WithOption("deployment", "gpt-4.1-nano"),
					WithOption("api_type", "azure"),
					WithOption("api_version", "2023-05-15"),
					WithOption("embedding_model", "text-embedding-ada-002"),
				),
			),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summarizer, err := NewLLMSummarizerFromConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}

				if summarizer != nil {
					t.Errorf("Expected nil summarizer when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if summarizer == nil {
					t.Errorf("Expected summarizer but got nil")
				}
			}
		})
	}
}

func TestLLMSummarizer_Summarize(t *testing.T) {
	// Skip if no API key is available
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("AZURE_OPENAI_API_KEY not set, skipping LLM tests")
	}

	// Test with Azure OpenAI
	config := NewConfig(
		WithType(TypeLlm),
		WithOpenAI("gpt-4", apiKey,
			WithURL(os.Getenv("AZURE_OPENAI_ENDPOINT")),
			WithOption("deployment", os.Getenv("AZURE_OPENAI_DEPLOYMENT")),
		),
	)

	summarizer, err := NewLLMSummarizerFromConfig(*config)
	if err != nil {
		t.Fatalf("Failed to create LLM summarizer: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expectError bool
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
			input:       "This is a single sentence to summarize.",
			expectError: false,
		},
		{
			name: "multiple sentences",
			input: `This is the first sentence of a longer text. 
					This is the second sentence that provides more context. 
					This is the third sentence that adds additional information. 
					This is the fourth sentence that continues the narrative. 
					This is the fifth sentence that concludes the text.`,
			expectError: false,
		},
		{
			name:        "German text",
			input:       "Dies ist ein deutscher Text, der zusammengefasst werden soll.",
			expectError: false,
		},
		{
			name: "Long text",
			input: `Artificial intelligence (AI) is intelligence demonstrated by machines, 
					in contrast to the natural intelligence displayed by humans and animals. 
					Leading AI textbooks define the field as the study of "intelligent agents": 
					any device that perceives its environment and takes actions that maximize 
					its chance of successfully achieving its goals. Colloquially, the term 
					"artificial intelligence" is often used to describe machines (or computers) 
					that mimic "cognitive" functions that humans associate with the human mind, 
					such as "learning" and "problem solving".`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := summarizer.Summarize(ctx, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if result == "" {
					t.Error("Expected non-empty result")
				}

				// For longer texts, the summary should be shorter
				if len(tt.input) > 100 && len(result) >= len(tt.input) {
					t.Errorf("Summary should be shorter than input: summary=%d, input=%d", len(result), len(tt.input))
				}

				t.Logf("Input: %s", tt.input)
				t.Logf("Summary: %s", result)
			}
		})
	}
}

func TestLLMSummarizer_WithContext(t *testing.T) {
	// Skip if no API key is available
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("AZURE_OPENAI_API_KEY not set, skipping context tests")
	}

	config := NewConfig(
		WithType(TypeLlm),
		WithOpenAI("gpt-4", apiKey,
			WithURL(os.Getenv("AZURE_OPENAI_ENDPOINT")),
			WithOption("deployment", os.Getenv("AZURE_OPENAI_DEPLOYMENT")),
		),
	)

	summarizer, err := NewLLMSummarizerFromConfig(*config)
	if err != nil {
		t.Fatalf("Failed to create LLM summarizer: %v", err)
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

// Add a test for NewLLMSummarizer with a dummy client
func TestNewLLMSummarizer_WithClient(t *testing.T) {
	dummy := &dummyLLMClient{}
	summarizer := NewLLMSummarizer(dummy)
	ctx := context.Background()

	result, err := summarizer.Summarize(ctx, "Test input")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "dummy-response" {
		t.Errorf("Expected dummy-response, got %s", result)
	}
}

type dummyLLMClient struct{}

func (d *dummyLLMClient) Generate(ctx context.Context, _ string) (string, error) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "dummy-response", nil
	}
}
