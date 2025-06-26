package summarizer

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Test default config
	config := NewConfig()
	expected := &Config{
		Type: TypeLexrank,
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Default config = %v, want %v", config, expected)
	}

	// Test with LLM config
	config2 := NewConfig(
		WithType(TypeLlm),
		WithLLM(
			WithProvider(LLMProviderOpenAI),
			WithModel("gpt-4"),
			WithAPIKey("test-key"),
		),
	)
	expected2 := &Config{
		Type: TypeLlm,
		LLM: &LLMConfig{
			Provider: LLMProviderOpenAI,
			Model:    "gpt-4",
			APIKey:   "test-key",
			Options:  make(map[string]interface{}),
		},
	}

	if !reflect.DeepEqual(config2, expected2) {
		t.Errorf("LLM config = %v, want %v", config2, expected2)
	}
}

func TestWithType(t *testing.T) {
	config := NewConfig()

	// Test default
	if config.Type != TypeLexrank {
		t.Errorf("Default type = %v, want %v", config.Type, TypeLexrank)
	}

	// Test setting type
	WithType(TypeLlm)(config)

	if config.Type != TypeLlm {
		t.Errorf("Type after WithType = %v, want %v", config.Type, TypeLlm)
	}
}

func TestWithLLM(t *testing.T) {
	config := NewConfig()

	// Test LLM configuration
	WithLLM(
		WithProvider(LLMProviderAnthropic),
		WithModel("claude-3-sonnet"),
		WithAPIKey("test-key"),
		WithMaxTokens(1000),
	)(config)

	if config.LLM == nil {
		t.Fatal("LLM should not be nil")
	}

	if config.LLM.Provider != LLMProviderAnthropic {
		t.Errorf("Provider = %v, want %v", config.LLM.Provider, LLMProviderAnthropic)
	}

	if config.LLM.Model != "claude-3-sonnet" {
		t.Errorf("Model = %v, want %v", config.LLM.Model, "claude-3-sonnet")
	}

	if config.LLM.APIKey != "test-key" {
		t.Errorf("APIKey = %v, want %v", config.LLM.APIKey, "test-key")
	}

	if config.LLM.MaxTokens != 1000 {
		t.Errorf("MaxTokens = %v, want %v", config.LLM.MaxTokens, 1000)
	}

	if config.LLM.Options == nil {
		t.Error("Options map should be initialized")
	}
}

func TestLLMOptions(t *testing.T) {
	llmConfig := &LLMConfig{Options: make(map[string]interface{})}

	// Test WithProvider
	WithProvider(LLMProviderMistral)(llmConfig)

	if llmConfig.Provider != LLMProviderMistral {
		t.Errorf("Provider = %v, want %v", llmConfig.Provider, LLMProviderMistral)
	}

	// Test WithModel
	WithModel("mistral-large")(llmConfig)

	if llmConfig.Model != "mistral-large" {
		t.Errorf("Model = %v, want %v", llmConfig.Model, "mistral-large")
	}

	// Test WithAPIKey
	WithAPIKey("sk-test")(llmConfig)

	if llmConfig.APIKey != "sk-test" {
		t.Errorf("APIKey = %v, want %v", llmConfig.APIKey, "sk-test")
	}

	// Test WithURL
	WithURL("https://api.mistral.ai/v1")(llmConfig)

	if llmConfig.URL != "https://api.mistral.ai/v1" {
		t.Errorf("URL = %v, want %v", llmConfig.URL, "https://api.mistral.ai/v1")
	}

	// Test WithBaseURL
	WithBaseURL("https://custom.anthropic.com")(llmConfig)

	if llmConfig.BaseURL != "https://custom.anthropic.com" {
		t.Errorf("BaseURL = %v, want %v", llmConfig.BaseURL, "https://custom.anthropic.com")
	}

	// Test WithMaxTokens
	WithMaxTokens(2000)(llmConfig)

	if llmConfig.MaxTokens != 2000 {
		t.Errorf("MaxTokens = %v, want %v", llmConfig.MaxTokens, 2000)
	}

	// Test WithAccountID
	WithAccountID("test-account")(llmConfig)

	if llmConfig.AccountID != "test-account" {
		t.Errorf("AccountID = %v, want %v", llmConfig.AccountID, "test-account")
	}

	// Test WithCredentials
	WithCredentials("/path/to/creds.json", "")(llmConfig)

	if llmConfig.Credentials == nil {
		t.Fatal("Credentials should not be nil")
	}

	if llmConfig.Credentials.Path != "/path/to/creds.json" {
		t.Errorf("Credentials.Path = %v, want %v", llmConfig.Credentials.Path, "/path/to/creds.json")
	}

	// Test WithOption
	WithOption("temperature", 0.7)(llmConfig)

	if llmConfig.Options["temperature"] != 0.7 {
		t.Errorf("Options[temperature] = %v, want %v", llmConfig.Options["temperature"], 0.7)
	}
}

func TestProviderConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name     string
		option   Option
		expected *LLMConfig
	}{
		{
			name:   "WithOpenAI",
			option: WithOpenAI("gpt-4", "sk-test", WithMaxTokens(1000)),
			expected: &LLMConfig{
				Provider:  LLMProviderOpenAI,
				Model:     "gpt-4",
				APIKey:    "sk-test",
				MaxTokens: 1000,
				Options:   make(map[string]interface{}),
			},
		},
		{
			name:   "WithAnthropic",
			option: WithAnthropic("claude-3-sonnet", "sk-ant-test"),
			expected: &LLMConfig{
				Provider: LLMProviderAnthropic,
				Model:    "claude-3-sonnet",
				APIKey:   "sk-ant-test",
				Options:  make(map[string]interface{}),
			},
		},
		{
			name:   "WithMistral",
			option: WithMistral("mistral-large", "sk-test", "https://api.mistral.ai/v1"),
			expected: &LLMConfig{
				Provider: LLMProviderMistral,
				Model:    "mistral-large",
				APIKey:   "sk-test",
				URL:      "https://api.mistral.ai/v1",
				Options:  make(map[string]interface{}),
			},
		},
		{
			name:   "WithGemini",
			option: WithGemini("gemini-pro", WithCredentials("/path/to/creds.json", "")),
			expected: &LLMConfig{
				Provider: LLMProviderGemini,
				Model:    "gemini-pro",
				Credentials: &Credentials{
					Path: "/path/to/creds.json",
					JSON: "",
				},
				Options: make(map[string]interface{}),
			},
		},
		{
			name:   "WithHuggingFace",
			option: WithHuggingFace("microsoft/DialoGPT-medium", "hf-test", "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium"),
			expected: &LLMConfig{
				Provider: LLMProviderHuggingFace,
				Model:    "microsoft/DialoGPT-medium",
				APIKey:   "hf-test",
				URL:      "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium",
				Options:  make(map[string]interface{}),
			},
		},
		{
			name:   "WithOllama",
			option: WithOllama("llama2", "http://localhost:11434"),
			expected: &LLMConfig{
				Provider: LLMProviderOllama,
				Model:    "llama2",
				URL:      "http://localhost:11434",
				Options:  make(map[string]interface{}),
			},
		},
		{
			name:   "WithCloudflare",
			option: WithCloudflare("@cf/meta/llama-2-7b-chat-int8", "cf-test", "test-account"),
			expected: &LLMConfig{
				Provider:  LLMProviderCloudflare,
				Model:     "@cf/meta/llama-2-7b-chat-int8",
				APIKey:    "cf-test",
				AccountID: "test-account",
				Options:   make(map[string]interface{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.option)

			if config.LLM == nil {
				t.Fatal("LLM should not be nil")
			}

			if config.LLM.Provider != tt.expected.Provider {
				t.Errorf("Provider = %v, want %v", config.LLM.Provider, tt.expected.Provider)
			}

			if config.LLM.Model != tt.expected.Model {
				t.Errorf("Model = %v, want %v", config.LLM.Model, tt.expected.Model)
			}

			if config.LLM.APIKey != tt.expected.APIKey {
				t.Errorf("APIKey = %v, want %v", config.LLM.APIKey, tt.expected.APIKey)
			}

			if config.LLM.URL != tt.expected.URL {
				t.Errorf("URL = %v, want %v", config.LLM.URL, tt.expected.URL)
			}

			if config.LLM.AccountID != tt.expected.AccountID {
				t.Errorf("AccountID = %v, want %v", config.LLM.AccountID, tt.expected.AccountID)
			}

			if config.LLM.MaxTokens != tt.expected.MaxTokens {
				t.Errorf("MaxTokens = %v, want %v", config.LLM.MaxTokens, tt.expected.MaxTokens)
			}

			if tt.expected.Credentials != nil {
				if config.LLM.Credentials == nil {
					t.Fatal("Credentials should not be nil")
				}

				if config.LLM.Credentials.Path != tt.expected.Credentials.Path {
					t.Errorf("Credentials.Path = %v, want %v", config.LLM.Credentials.Path, tt.expected.Credentials.Path)
				}
			}
		})
	}
}

func TestComplexConfiguration(t *testing.T) {
	// Test a complex configuration with multiple options
	config := NewConfig(
		WithType(TypeLlm),
		WithOpenAI("gpt-4", "sk-test",
			WithMaxTokens(1500),
			WithURL("https://api.openai.com/v1"),
			WithOption("temperature", 0.8),
			WithOption("organization_id", "org-test"),
		),
	)

	if config.Type != TypeLlm {
		t.Errorf("Type = %v, want %v", config.Type, TypeLlm)
	}

	if config.LLM == nil {
		t.Fatal("LLM should not be nil")
	}

	if config.LLM.Provider != LLMProviderOpenAI {
		t.Errorf("Provider = %v, want %v", config.LLM.Provider, LLMProviderOpenAI)
	}

	if config.LLM.Model != "gpt-4" {
		t.Errorf("Model = %v, want %v", config.LLM.Model, "gpt-4")
	}

	if config.LLM.APIKey != "sk-test" {
		t.Errorf("APIKey = %v, want %v", config.LLM.APIKey, "sk-test")
	}

	if config.LLM.MaxTokens != 1500 {
		t.Errorf("MaxTokens = %v, want %v", config.LLM.MaxTokens, 1500)
	}

	if config.LLM.URL != "https://api.openai.com/v1" {
		t.Errorf("URL = %v, want %v", config.LLM.URL, "https://api.openai.com/v1")
	}

	if config.LLM.Options["temperature"] != 0.8 {
		t.Errorf("Options[temperature] = %v, want %v", config.LLM.Options["temperature"], 0.8)
	}

	if config.LLM.Options["organization_id"] != "org-test" {
		t.Errorf("Options[organization_id] = %v, want %v", config.LLM.Options["organization_id"], "org-test")
	}
}

func TestOptionsMapInitialization(t *testing.T) {
	// Test that Options map is properly initialized
	config := NewConfig(
		WithLLM(
			WithProvider(LLMProviderOpenAI),
			WithModel("gpt-4"),
		),
	)

	if config.LLM.Options == nil {
		t.Fatal("Options map should be initialized")
	}

	// Test adding options
	config.LLM.Options["test"] = "value"
	if config.LLM.Options["test"] != "value" {
		t.Errorf("Options[test] = %v, want %v", config.LLM.Options["test"], "value")
	}
}

func TestCredentialsHandling(t *testing.T) {
	// Test credentials with path only
	config1 := NewConfig(
		WithLLM(
			WithProvider(LLMProviderGemini),
			WithModel("gemini-pro"),
			WithCredentials("/path/to/creds.json", ""),
		),
	)

	if config1.LLM.Credentials.Path != "/path/to/creds.json" {
		t.Errorf("Credentials.Path = %v, want %v", config1.LLM.Credentials.Path, "/path/to/creds.json")
	}

	// Test credentials with JSON only
	config2 := NewConfig(
		WithLLM(
			WithProvider(LLMProviderGemini),
			WithModel("gemini-pro"),
			WithCredentials("", `{"test": "credentials"}`),
		),
	)

	if config2.LLM.Credentials.JSON != `{"test": "credentials"}` {
		t.Errorf("Credentials.JSON = %v, want %v", config2.LLM.Credentials.JSON, `{"test": "credentials"}`)
	}

	// Test credentials with both path and JSON
	config3 := NewConfig(
		WithLLM(
			WithProvider(LLMProviderGemini),
			WithModel("gemini-pro"),
			WithCredentials("/path/to/creds.json", `{"test": "credentials"}`),
		),
	)

	if config3.LLM.Credentials.Path != "/path/to/creds.json" {
		t.Errorf("Credentials.Path = %v, want %v", config3.LLM.Credentials.Path, "/path/to/creds.json")
	}

	if config3.LLM.Credentials.JSON != `{"test": "credentials"}` {
		t.Errorf("Credentials.JSON = %v, want %v", config3.LLM.Credentials.JSON, `{"test": "credentials"}`)
	}
}

func TestSummarizerTypeConstants(t *testing.T) {
	// Test that constants have expected values
	if TypeLexrank != "lexrank" {
		t.Errorf("TypeLexrank = %v, want %v", TypeLexrank, "lexrank")
	}

	if TypeLlm != "llm" {
		t.Errorf("TypeLlm = %v, want %v", TypeLlm, "llm")
	}
}

func TestLLMProviderConstants(t *testing.T) {
	// Test that all provider constants have expected values
	expectedProviders := map[LLMProvider]string{
		LLMProviderOpenAI:      "openai",
		LLMProviderAnthropic:   "anthropic",
		LLMProviderMistral:     "mistral",
		LLMProviderGemini:      "gemini",
		LLMProviderCloudflare:  "cloudflare",
		LLMProviderHuggingFace: "huggingface",
		LLMProviderOllama:      "ollama",
	}

	for provider, expected := range expectedProviders {
		if string(provider) != expected {
			t.Errorf("Provider %v = %v, want %v", provider, string(provider), expected)
		}
	}
}
