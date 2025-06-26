package llm

// Provider represents the supported LLM service providers.
//
// Each provider has different API endpoints, authentication methods, and
// configuration requirements. The provider determines which fields in
// Config are required and how they are used.
type Provider string

const (
	// ProviderOpenAI represents OpenAI's API service.
	// Requires: Model, APIKey
	// Optional: URL, MaxTokens, Options["organization_id"]
	ProviderOpenAI Provider = "openai"

	// ProviderAnthropic represents Anthropic's Claude API service.
	// Requires: Model, APIKey
	// Optional: BaseURL, MaxTokens, Options["beta_header"]
	ProviderAnthropic Provider = "anthropic"

	// ProviderMistral represents Mistral AI's API service.
	// Requires: Model, APIKey, URL
	// Optional: MaxTokens
	ProviderMistral Provider = "mistral"

	// ProviderGemini represents Google's Gemini API service.
	// Requires: Model
	// Optional: Credentials, MaxTokens
	ProviderGemini Provider = "gemini"

	// ProviderCloudflare represents Cloudflare's AI service.
	// Requires: Model, APIKey, AccountID
	// Optional: URL, MaxTokens
	ProviderCloudflare Provider = "cloudflare"

	// ProviderHuggingFace represents HuggingFace's inference API.
	// Requires: Model, APIKey, URL
	// Optional: MaxTokens
	ProviderHuggingFace Provider = "huggingface"

	// ProviderOllama represents local Ollama deployments.
	// Requires: Model, URL
	// Optional: MaxTokens
	ProviderOllama Provider = "ollama"
)

// Config contains all configuration parameters for LLM-based services.
//
// This struct consolidates configuration for all supported LLM providers into
// a single, unified structure. Provider-specific settings are stored in the
// Options map for extensibility.
type Config struct {
	// Provider specifies which LLM service to use.
	// Must be one of the predefined Provider constants.
	Provider Provider

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

// Option is a function that modifies a Config instance.
//
// Options are used with NewConfig to create configured instances.
// This pattern provides a flexible, readable way to configure the LLM client.
type Option func(*Config)

// NewConfig creates a new Config instance with the specified options.
//
// The function applies each option in sequence to build the final configuration.
// If no options are provided, an empty configuration is returned.
func NewConfig(options ...Option) *Config {
	config := &Config{
		Options: make(map[string]interface{}),
	}

	for _, option := range options {
		option(config)
	}

	return config
}

// WithProvider sets the LLM service provider.
//
// This option must be called when configuring LLM-based services.
// The provider determines which other configuration options are required
// and how they are interpreted.
func WithProvider(provider Provider) Option {
	return func(c *Config) {
		c.Provider = provider
	}
}

// WithModel sets the model name for the LLM service.
//
// The model name is provider-specific and determines the capabilities
// and pricing of the service.
//
// Examples:
//   - OpenAI: "gpt-4", "gpt-3.5-turbo"
//   - Anthropic: "claude-3-sonnet", "claude-3-haiku"
//   - Mistral: "mistral-large", "mistral-medium"
func WithModel(model string) Option {
	return func(c *Config) {
		c.Model = model
	}
}

// WithAPIKey sets the authentication key for the LLM service.
//
// Most LLM providers require an API key for authentication. The key format
// and requirements vary by provider. Local services like Ollama typically
// do not require an API key.
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// WithURL sets the API endpoint URL for the LLM service.
//
// This option is used for providers that support custom endpoints or
// local deployments. For example, it can be used to point to a local
// Ollama instance or a custom Mistral deployment.
func WithURL(url string) Option {
	return func(c *Config) {
		c.URL = url
	}
}

// WithBaseURL sets the base URL for the LLM service.
//
// This option is used by providers like Anthropic that support custom
// base URLs for different regions or deployments.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithMaxTokens sets the maximum number of tokens in the response.
//
// This option helps control response length and can reduce API costs.
// The actual token limit depends on the specific model being used.
func WithMaxTokens(maxTokens int) Option {
	return func(c *Config) {
		c.MaxTokens = maxTokens
	}
}

// WithAccountID sets the account identifier for the LLM service.
//
// This option is required for some providers like Cloudflare that use
// account-based authentication and resource management.
func WithAccountID(accountID string) Option {
	return func(c *Config) {
		c.AccountID = accountID
	}
}

// WithCredentials sets the authentication credentials for the LLM service.
//
// This option is used for providers that require file-based or JSON-based
// authentication, such as Google Cloud services.
//
// Either path or json should be provided, depending on the credential format.
func WithCredentials(path, json string) Option {
	return func(c *Config) {
		c.Credentials = &Credentials{
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
func WithOption(key string, value interface{}) Option {
	return func(c *Config) {
		c.Options[key] = value
	}
}

// Convenience functions for specific providers

// WithOpenAI creates a complete OpenAI configuration.
//
// This function provides a convenient way to configure OpenAI services
// with the most commonly used parameters. Additional options can be provided
// for advanced configuration.
//
// Example:
//
//	config := NewConfig(
//		WithOpenAI("gpt-4", "sk-...", WithMaxTokens(1000)),
//	)
func WithOpenAI(model, apiKey string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderOpenAI
		c.Model = model
		c.APIKey = apiKey

		for _, option := range options {
			option(c)
		}
	}
}

// WithAnthropic creates a complete Anthropic configuration.
//
// This function provides a convenient way to configure Anthropic Claude
// services with the most commonly used parameters.
//
// Example:
//
//	config := NewConfig(
//		WithAnthropic("claude-3-sonnet", "sk-ant-...", WithMaxTokens(1000)),
//	)
func WithAnthropic(model, apiKey string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderAnthropic
		c.Model = model
		c.APIKey = apiKey

		for _, option := range options {
			option(c)
		}
	}
}

// WithMistral creates a complete Mistral configuration.
//
// This function provides a convenient way to configure Mistral AI
// services. Mistral requires a URL parameter for the API endpoint.
//
// Example:
//
//	config := NewConfig(
//		WithMistral("mistral-large", "sk-...", "https://api.mistral.ai/v1"),
//	)
func WithMistral(model, apiKey, url string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderMistral
		c.Model = model
		c.APIKey = apiKey
		c.URL = url

		for _, option := range options {
			option(c)
		}
	}
}

// WithGemini creates a complete Google Gemini configuration.
//
// This function provides a convenient way to configure Google Gemini
// services. Gemini typically uses service account credentials
// rather than API keys.
//
// Example:
//
//	config := NewConfig(
//		WithGemini("gemini-pro", WithCredentials("/path/to/credentials.json", "")),
//	)
func WithGemini(model string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderGemini
		c.Model = model

		for _, option := range options {
			option(c)
		}
	}
}

// WithHuggingFace creates a complete HuggingFace configuration.
//
// This function provides a convenient way to configure HuggingFace
// inference API services. HuggingFace requires a URL parameter
// for the specific model endpoint.
//
// Example:
//
//	config := NewConfig(
//		WithHuggingFace("microsoft/DialoGPT-medium", "hf_...", "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium"),
//	)
func WithHuggingFace(model, apiKey, url string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderHuggingFace
		c.Model = model
		c.APIKey = apiKey
		c.URL = url

		for _, option := range options {
			option(c)
		}
	}
}

// WithOllama creates a complete Ollama configuration.
//
// This function provides a convenient way to configure local Ollama
// services. Ollama typically runs locally and doesn't require
// an API key.
//
// Example:
//
//	config := NewConfig(
//		WithOllama("llama2", "http://localhost:11434"),
//	)
func WithOllama(model, url string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderOllama
		c.Model = model
		c.URL = url

		for _, option := range options {
			option(c)
		}
	}
}

// WithCloudflare creates a complete Cloudflare configuration.
//
// This function provides a convenient way to configure Cloudflare AI
// services. Cloudflare requires an account ID for authentication.
//
// Example:
//
//	config := NewConfig(
//		WithCloudflare("@cf/meta/llama-2-7b-chat-int8", "cf_...", "your-account-id"),
//	)
func WithCloudflare(model, apiKey, accountID string, options ...Option) Option {
	return func(c *Config) {
		c.Provider = ProviderCloudflare
		c.Model = model
		c.APIKey = apiKey
		c.AccountID = accountID

		for _, option := range options {
			option(c)
		}
	}
}
