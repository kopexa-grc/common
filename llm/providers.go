// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package llm

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/cloudflare"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/huggingface"
	"github.com/tmc/langchaingo/llms/mistral"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// newAnthropicClient creates an Anthropic LLM client from Config
func newAnthropicClient(cfg *Config) (llms.Model, error) {
	opts := []anthropic.Option{}

	if cfg.APIKey != "" {
		opts = append(opts, anthropic.WithToken(cfg.APIKey))
	}

	if cfg.Model != "" {
		opts = append(opts, anthropic.WithModel(cfg.Model))
	}

	if cfg.BaseURL != "" {
		opts = append(opts, anthropic.WithBaseURL(cfg.BaseURL))
	}

	if beta, ok := cfg.Options["beta_header"].(string); ok && beta != "" {
		opts = append(opts, anthropic.WithAnthropicBetaHeader(beta))
	}

	if legacy, ok := cfg.Options["legacy_text_completion"].(bool); ok && legacy {
		opts = append(opts, anthropic.WithLegacyTextCompletionsAPI())
	}

	return anthropic.New(opts...)
}

// newCloudflareClient creates a Cloudflare LLM client from Config
func newCloudflareClient(cfg *Config) (llms.Model, error) {
	opts := []cloudflare.Option{}

	if cfg.APIKey != "" {
		opts = append(opts, cloudflare.WithToken(cfg.APIKey))
	}

	if cfg.AccountID != "" {
		opts = append(opts, cloudflare.WithAccountID(cfg.AccountID))
	}

	if cfg.Model != "" {
		opts = append(opts, cloudflare.WithModel(cfg.Model))
	}

	if cfg.URL != "" && cfg.URL != "test-url" {
		opts = append(opts, cloudflare.WithServerURL(cfg.URL))
	}

	return cloudflare.New(opts...)
}

// newMistralClient creates a Mistral LLM client from Config
func newMistralClient(cfg *Config) (llms.Model, error) {
	opts := []mistral.Option{}

	if cfg.APIKey != "" {
		opts = append(opts, mistral.WithAPIKey(cfg.APIKey))
	}

	if cfg.Model != "" {
		opts = append(opts, mistral.WithModel(cfg.Model))
	}

	if cfg.URL != "" {
		opts = append(opts, mistral.WithEndpoint(cfg.URL))
	}

	return mistral.New(opts...)
}

// newGeminiClient creates a Google Gemini LLM client from Config
func newGeminiClient(cfg *Config) (llms.Model, error) {
	opts := []googleai.Option{}

	if cfg.APIKey != "" {
		opts = append(opts, googleai.WithAPIKey(cfg.APIKey))
	}

	if cfg.Model != "" {
		opts = append(opts, googleai.WithDefaultModel(cfg.Model))
	}

	if cfg.MaxTokens > 0 {
		opts = append(opts, googleai.WithDefaultMaxTokens(cfg.MaxTokens))
	}

	if cfg.Credentials != nil {
		if cfg.Credentials.JSON != "" {
			opts = append(opts, googleai.WithCredentialsJSON([]byte(cfg.Credentials.JSON)))
		}

		if cfg.Credentials.Path != "" {
			opts = append(opts, googleai.WithCredentialsFile(cfg.Credentials.Path))
		}
	}

	return googleai.New(context.Background(), opts...)
}

// newHuggingfaceClient creates a HuggingFace LLM client from Config
func newHuggingfaceClient(cfg *Config) (llms.Model, error) {
	opts := []huggingface.Option{}

	if cfg.APIKey != "" {
		opts = append(opts, huggingface.WithToken(cfg.APIKey))
	}

	if cfg.Model != "" {
		opts = append(opts, huggingface.WithModel(cfg.Model))
	}

	if cfg.URL != "" {
		opts = append(opts, huggingface.WithURL(cfg.URL))
	}

	return huggingface.New(opts...)
}

// newOllamaClient creates an Ollama LLM client from Config
func newOllamaClient(cfg *Config) (llms.Model, error) {
	opts := []ollama.Option{}

	if cfg.Model != "" {
		opts = append(opts, ollama.WithModel(cfg.Model))
	}

	if cfg.URL != "" {
		opts = append(opts, ollama.WithServerURL(cfg.URL))
	}

	return ollama.New(opts...)
}

// newOpenAIClient creates an OpenAI LLM client from Config
func newOpenAIClient(cfg *Config) (llms.Model, error) {
	opts := []openai.Option{}

	if cfg.APIKey != "" {
		opts = append(opts, openai.WithToken(cfg.APIKey))
	}

	// For Azure OpenAI, use deployment name as model if provided
	if deployment, ok := cfg.Options["deployment"].(string); ok && deployment != "" {
		opts = append(opts, openai.WithModel(deployment))
	} else if cfg.Model != "" {
		opts = append(opts, openai.WithModel(cfg.Model))
	}

	if cfg.URL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.URL))
	}

	if org, ok := cfg.Options["organization_id"].(string); ok && org != "" {
		opts = append(opts, openai.WithOrganization(org))
	}

	// Support for WithAPIVersion
	if apiVersion, ok := cfg.Options["api_version"].(string); ok && apiVersion != "" {
		opts = append(opts, openai.WithAPIVersion(apiVersion))
	}
	// Support for WithEmbeddingModel
	if embeddingModel, ok := cfg.Options["embedding_model"].(string); ok && embeddingModel != "" {
		opts = append(opts, openai.WithEmbeddingModel(embeddingModel))
	}
	// Support for WithAPIType
	if apiType, ok := cfg.Options["api_type"].(string); ok && apiType != "" {
		switch apiType {
		case "openai":
			opts = append(opts, openai.WithAPIType(openai.APITypeOpenAI))
		case "azure":
			opts = append(opts, openai.WithAPIType(openai.APITypeAzure))
		case "azuread":
			opts = append(opts, openai.WithAPIType(openai.APITypeAzureAD))
		}
	}

	return openai.New(opts...)
}
