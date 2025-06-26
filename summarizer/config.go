// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package summarizer provides configuration management for text summarization services.
//
// This package implements the Options Pattern for flexible and type-safe configuration
// of different summarization backends, including LexRank and various LLM providers.
//
// Example usage:
//
//	config := summarizer.NewConfig(
//		summarizer.WithType(summarizer.TypeLlm),
//		summarizer.WithOpenAI("gpt-4", "your-api-key"),
//	)
//
// The package supports multiple LLM providers including OpenAI, Anthropic, Mistral,
// Google Gemini, HuggingFace, Ollama, and Cloudflare.
package summarizer

// Type represents the type of summarization algorithm to use.
type Type string

const (
	// TypeLexrank uses the LexRank algorithm for extractive summarization.
	// This method ranks sentences based on their centrality in the document graph.
	TypeLexrank Type = "lexrank"

	// TypeLlm uses a Large Language Model for abstractive summarization.
	// This method generates new text that captures the key information from the source.
	TypeLlm Type = "llm"
)

// Config represents the complete configuration for a summarization service.
//
// The configuration supports both extractive (LexRank) and abstractive (LLM)
// summarization methods. When using LLM-based summarization, additional LLM
// configuration must be provided.
type Config struct {
	// Type specifies the summarization algorithm to use.
	// Defaults to TypeLexrank if not specified.
	Type Type

	// LLM contains the configuration for LLM-based summarization.
	// Required when Type is TypeLlm, ignored otherwise.
	LLM *LLMConfig
}

// LLMConfig contains all configuration parameters for LLM-based summarization.
//
// This struct consolidates configuration for all supported LLM providers into
// a single, unified structure. Provider-specific settings are stored in the
// Options map for extensibility.
type LLMConfig struct {
	// Provider specifies which LLM service to use.
	// Must be one of the predefined LLMProvider constants.
	Provider LLMProvider

	// Model specifies the specific model name to use with the provider.
	// Examples: "gpt-4", "claude-3-sonnet", "llama2"
	Model string

	// APIKey contains the authentication key for the LLM service.
	// Required for most providers except local services like Ollama.
	APIKey string

	// URL specifies the API endpoint URL.
	// Used for providers that support custom endpoints or local deployments.
	URL string

	// BaseURL specifies the base URL for the API service.
	// Used by providers like Anthropic that support custom base URLs.
	BaseURL string

	// MaxTokens specifies the maximum number of tokens in the response.
	// Helps control response length and API costs.
	MaxTokens int

	// AccountID specifies the account identifier for the service.
	// Required for some providers like Cloudflare.
	AccountID string

	// Credentials contains authentication credentials for the service.
	// Used for providers that require file-based or JSON credentials.
	Credentials *Credentials

	// Options contains provider-specific configuration parameters.
	// This map allows for extensible configuration without struct changes.
	// Common keys include "temperature", "organization_id", "beta_header".
	Options map[string]interface{}
}

// Credentials represents authentication credentials for LLM services.
//
// This struct supports both file-based and JSON-based credential storage
// for services that require more complex authentication methods.
type Credentials struct {
	// Path specifies the file path to the credentials file.
	// Used for Google Cloud and other file-based authentication systems.
	Path string

	// JSON contains the credentials as a JSON string.
	// Used when credentials need to be provided as a string rather than a file.
	JSON string
}

// LLMProvider represents the supported LLM service providers.
//
// Each provider has different API endpoints, authentication methods, and
// configuration requirements. The provider determines which fields in
// LLMConfig are required and how they are used.
type LLMProvider string

const (
	// LLMProviderOpenAI represents OpenAI's API service.
	// Requires: Model, APIKey
	// Optional: URL, MaxTokens, Options["organization_id"]
	LLMProviderOpenAI LLMProvider = "openai"

	// LLMProviderAnthropic represents Anthropic's Claude API service.
	// Requires: Model, APIKey
	// Optional: BaseURL, MaxTokens, Options["beta_header"]
	LLMProviderAnthropic LLMProvider = "anthropic"

	// LLMProviderMistral represents Mistral AI's API service.
	// Requires: Model, APIKey, URL
	// Optional: MaxTokens
	LLMProviderMistral LLMProvider = "mistral"

	// LLMProviderGemini represents Google's Gemini API service.
	// Requires: Model
	// Optional: Credentials, MaxTokens
	LLMProviderGemini LLMProvider = "gemini"

	// LLMProviderCloudflare represents Cloudflare's AI service.
	// Requires: Model, APIKey, AccountID
	// Optional: URL, MaxTokens
	LLMProviderCloudflare LLMProvider = "cloudflare"

	// LLMProviderHuggingFace represents HuggingFace's inference API.
	// Requires: Model, APIKey, URL
	// Optional: MaxTokens
	LLMProviderHuggingFace LLMProvider = "huggingface"

	// LLMProviderOllama represents local Ollama deployments.
	// Requires: Model, URL
	// Optional: MaxTokens
	LLMProviderOllama LLMProvider = "ollama"
)

// Option is a function that modifies a Config instance.
//
// Options are used with NewConfig to create configured instances.
// This pattern provides a flexible, readable way to configure the summarizer.
type Option func(*Config)

// LLMOption is a function that modifies an LLMConfig instance.
//
// LLMOptions are used with WithLLM to configure LLM-specific settings.
// This allows for granular control over LLM configuration.
type LLMOption func(*LLMConfig)

// NewConfig creates a new Config instance with the specified options.
//
// The function applies each option in sequence to build the final configuration.
// If no options are provided, a default configuration with LexRank summarization
// is returned.
//
// Example:
//
//	config := NewConfig(
//		WithType(SummarizerTypeLlm),
//		WithOpenAI("gpt-4", "sk-..."),
//	)
func NewConfig(options ...Option) *Config {
	config := &Config{
		Type: TypeLexrank, // default to LexRank
	}

	for _, option := range options {
		option(config)
	}

	return config
}

// WithType sets the summarization type for the configuration.
//
// This option determines whether to use extractive (LexRank) or abstractive (LLM)
// summarization. When using LLM summarization, additional LLM configuration
// must be provided using WithLLM or one of the provider-specific functions.
func WithType(summarizerType Type) Option {
	return func(c *Config) {
		c.Type = summarizerType
	}
}

// WithLLM configures LLM-based summarization with the specified options.
//
// This option sets up the LLM configuration and should be used when Type is
// TypeLlm. The function creates a new LLMConfig and applies all
// provided LLMOptions to it.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithLLM(
//			WithProvider(LLMProviderOpenAI),
//			WithModel("gpt-4"),
//			WithAPIKey("sk-..."),
//		),
//	)
func WithLLM(options ...LLMOption) Option {
	return func(c *Config) {
		llmConfig := &LLMConfig{
			Options: make(map[string]interface{}),
		}

		for _, option := range options {
			option(llmConfig)
		}

		c.LLM = llmConfig
	}
}

// WithProvider sets the LLM service provider.
//
// This option must be called when configuring LLM-based summarization.
// The provider determines which other configuration options are required
// and how they are interpreted.
func WithProvider(provider LLMProvider) LLMOption {
	return func(l *LLMConfig) {
		l.Provider = provider
	}
}

// WithModel sets the model name for the LLM service.
//
// The model name is provider-specific and determines the capabilities
// and pricing of the summarization service.
//
// Examples:
//   - OpenAI: "gpt-4", "gpt-3.5-turbo"
//   - Anthropic: "claude-3-sonnet", "claude-3-haiku"
//   - Mistral: "mistral-large", "mistral-medium"
func WithModel(model string) LLMOption {
	return func(l *LLMConfig) {
		l.Model = model
	}
}

// WithAPIKey sets the authentication key for the LLM service.
//
// Most LLM providers require an API key for authentication. The key format
// and requirements vary by provider. Local services like Ollama typically
// do not require an API key.
func WithAPIKey(apiKey string) LLMOption {
	return func(l *LLMConfig) {
		l.APIKey = apiKey
	}
}

// WithURL sets the API endpoint URL for the LLM service.
//
// This option is used for providers that support custom endpoints or
// local deployments. For example, it can be used to point to a local
// Ollama instance or a custom Mistral deployment.
func WithURL(url string) LLMOption {
	return func(l *LLMConfig) {
		l.URL = url
	}
}

// WithBaseURL sets the base URL for the LLM service.
//
// This option is used by providers like Anthropic that support custom
// base URLs for different regions or deployments.
func WithBaseURL(baseURL string) LLMOption {
	return func(l *LLMConfig) {
		l.BaseURL = baseURL
	}
}

// WithMaxTokens sets the maximum number of tokens in the response.
//
// This option helps control response length and can reduce API costs.
// The actual token limit depends on the specific model being used.
func WithMaxTokens(maxTokens int) LLMOption {
	return func(l *LLMConfig) {
		l.MaxTokens = maxTokens
	}
}

// WithAccountID sets the account identifier for the LLM service.
//
// This option is required for some providers like Cloudflare that use
// account-based authentication and resource management.
func WithAccountID(accountID string) LLMOption {
	return func(l *LLMConfig) {
		l.AccountID = accountID
	}
}

// WithCredentials sets the authentication credentials for the LLM service.
//
// This option is used for providers that require file-based or JSON-based
// authentication, such as Google Cloud services.
//
// Either path or json should be provided, depending on the credential format.
func WithCredentials(path, json string) LLMOption {
	return func(l *LLMConfig) {
		l.Credentials = &Credentials{
			Path: path,
			JSON: json,
		}
	}
}

// WithOption sets a custom configuration option for the LLM service.
//
// This option allows for provider-specific configuration without requiring
// struct changes. Common options include:
//   - "temperature": Controls response randomness (0.0-1.0)
//   - "organization_id": OpenAI organization identifier
//   - "beta_header": Anthropic beta feature flags
//
// The interpretation of these options depends on the specific provider.
func WithOption(key string, value interface{}) LLMOption {
	return func(l *LLMConfig) {
		l.Options[key] = value
	}
}

// Convenience functions for specific providers

// WithOpenAI creates a complete OpenAI configuration.
//
// This function provides a convenient way to configure OpenAI summarization
// with the most commonly used parameters. Additional options can be provided
// for advanced configuration.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithOpenAI("gpt-4", "sk-...", WithMaxTokens(1000)),
//	)
func WithOpenAI(model, apiKey string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderOpenAI),
		WithModel(model),
		WithAPIKey(apiKey),
	}, options...)...)
}

// WithAnthropic creates a complete Anthropic configuration.
//
// This function provides a convenient way to configure Anthropic Claude
// summarization with the most commonly used parameters.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithAnthropic("claude-3-sonnet", "sk-ant-...", WithMaxTokens(1000)),
//	)
func WithAnthropic(model, apiKey string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderAnthropic),
		WithModel(model),
		WithAPIKey(apiKey),
	}, options...)...)
}

// WithMistral creates a complete Mistral configuration.
//
// This function provides a convenient way to configure Mistral AI
// summarization. Mistral requires a URL parameter for the API endpoint.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithMistral("mistral-large", "sk-...", "https://api.mistral.ai/v1"),
//	)
func WithMistral(model, apiKey, url string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderMistral),
		WithModel(model),
		WithAPIKey(apiKey),
		WithURL(url),
	}, options...)...)
}

// WithGemini creates a complete Google Gemini configuration.
//
// This function provides a convenient way to configure Google Gemini
// summarization. Gemini typically uses service account credentials
// rather than API keys.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithGemini("gemini-pro", WithCredentials("/path/to/credentials.json", "")),
//	)
func WithGemini(model string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderGemini),
		WithModel(model),
	}, options...)...)
}

// WithHuggingFace creates a complete HuggingFace configuration.
//
// This function provides a convenient way to configure HuggingFace
// inference API summarization. HuggingFace requires a URL parameter
// for the specific model endpoint.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithHuggingFace("microsoft/DialoGPT-medium", "hf_...", "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium"),
//	)
func WithHuggingFace(model, apiKey, url string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderHuggingFace),
		WithModel(model),
		WithAPIKey(apiKey),
		WithURL(url),
	}, options...)...)
}

// WithOllama creates a complete Ollama configuration.
//
// This function provides a convenient way to configure local Ollama
// summarization. Ollama typically runs locally and doesn't require
// an API key.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithOllama("llama2", "http://localhost:11434"),
//	)
func WithOllama(model, url string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderOllama),
		WithModel(model),
		WithURL(url),
	}, options...)...)
}

// WithCloudflare creates a complete Cloudflare configuration.
//
// This function provides a convenient way to configure Cloudflare AI
// summarization. Cloudflare requires an account ID for authentication.
//
// Example:
//
//	config := NewConfig(
//		WithType(TypeLlm),
//		WithCloudflare("@cf/meta/llama-2-7b-chat-int8", "cf_...", "your-account-id"),
//	)
func WithCloudflare(model, apiKey, accountID string, options ...LLMOption) Option {
	return WithLLM(append([]LLMOption{
		WithProvider(LLMProviderCloudflare),
		WithModel(model),
		WithAPIKey(apiKey),
		WithAccountID(accountID),
	}, options...)...)
}
