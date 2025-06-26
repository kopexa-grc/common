// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

// Client represents an LLM client that can be used for various text generation tasks.
type Client struct {
	llmClient llms.Model
}

// New creates a new LLM client with the given configuration.
//
// The function creates the appropriate LLM client based on the provider
// specified in the configuration. All supported providers are handled
// automatically.
//
// Example:
//
//	client, err := llm.New(llm.NewConfig(
//		llm.WithOpenAI("gpt-4", "sk-..."),
//	))
func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrConfigRequired
	}

	llmClient, err := getClient(*cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		llmClient: llmClient,
	}, nil
}

// Generate generates text based on the provided prompt.
//
// This method sends the prompt to the configured LLM and returns the generated response.
// The context can be used for cancellation and timeouts.
//
// Example:
//
//	result, err := client.Generate(ctx, "Summarize this text: ...")
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, c.llmClient, prompt)
}

// GenerateWithOptions generates text with additional options.
//
// This method allows for more control over the generation process by accepting
// additional options that are passed to the underlying LLM.
func (c *Client) GenerateWithOptions(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, c.llmClient, prompt, options...)
}

// GetModel returns the underlying LLM model for advanced usage.
//
// This method provides access to the underlying langchaingo model for cases
// where more advanced functionality is needed.
func (c *Client) GetModel() llms.Model {
	return c.llmClient
}

// getClient creates the appropriate LLM client based on the configuration
func getClient(cfg Config) (llms.Model, error) {
	switch cfg.Provider {
	case ProviderAnthropic:
		return newAnthropicClient(&cfg)
	case ProviderCloudflare:
		return newCloudflareClient(&cfg)
	case ProviderMistral:
		return newMistralClient(&cfg)
	case ProviderGemini:
		return newGeminiClient(&cfg)
	case ProviderHuggingFace:
		return newHuggingfaceClient(&cfg)
	case ProviderOllama:
		return newOllamaClient(&cfg)
	case ProviderOpenAI:
		return newOpenAIClient(&cfg)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, cfg.Provider)
	}
}
