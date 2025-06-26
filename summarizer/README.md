# Summarizer Package

The `summarizer` package provides a unified interface for text summarization using both extractive (LexRank) and abstractive (LLM-based) methods. It implements the Options Pattern for flexible configuration and supports multiple LLM providers through the `llm` package.

## Features

- **Multi-Method Support**: LexRank (extractive) and LLM-based (abstractive) summarization
- **Multi-Provider LLM Support**: OpenAI, Anthropic, Mistral, Google Gemini, HuggingFace, Ollama, Cloudflare
- **Azure OpenAI Support**: Complete Azure OpenAI integration with deployment management
- **Language Detection**: Automatic language detection with appropriate prompts
- **Input/Output Sanitization**: Built-in HTML sanitization for security
- **Flexible Configuration**: Options Pattern for type-safe configuration
- **Context Support**: Full context cancellation and timeout support
- **Testable Architecture**: Interface-based design for easy mocking and testing

## Installation

```bash
go get github.com/kopexa-grc/common/summarizer
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kopexa-grc/common/summarizer"
)

func main() {
    // Create a LexRank summarizer (default)
    client, err := summarizer.New(summarizer.NewConfig())
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    text := `This is a long text that needs to be summarized. 
             It contains multiple sentences and should be shortened 
             to capture the key points.`

    summary, err := client.Summarize(ctx, text)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(summary)
}
```

### LLM-Based Summarization

```go
// OpenAI
client, err := summarizer.New(summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLlm),
    summarizer.WithOpenAI("gpt-4", "your-api-key"),
))

// Azure OpenAI
client, err := summarizer.New(summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLlm),
    summarizer.WithOpenAI("gpt-4", "your-api-key",
        summarizer.WithURL("https://your-resource.openai.azure.com/"),
        summarizer.WithOption("deployment", "your-deployment-name"),
        summarizer.WithOption("api_type", "azure"),
        summarizer.WithOption("api_version", "2023-05-15"),
    ),
))

// Anthropic Claude
client, err := summarizer.New(summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLlm),
    summarizer.WithAnthropic("claude-3-sonnet", "your-api-key"),
))

// Local Ollama
client, err := summarizer.New(summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLlm),
    summarizer.WithOllama("llama2", "http://localhost:11434"),
))
```

## Configuration

### Options Pattern

The package uses the Options Pattern for flexible and type-safe configuration:

```go
config := summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLlm),
    summarizer.WithOpenAI("gpt-4", "your-api-key",
        summarizer.WithMaxTokens(1000),
        summarizer.WithOption("temperature", 0.7),
    ),
)
```

### Summarization Types

```go
const (
    TypeLexrank Type = "lexrank"  // Extractive summarization
    TypeLlm     Type = "llm"      // Abstractive summarization
)
```

### LLM Providers

```go
const (
    LLMProviderOpenAI      LLMProvider = "openai"
    LLMProviderAnthropic   LLMProvider = "anthropic"
    LLMProviderMistral     LLMProvider = "mistral"
    LLMProviderGemini      LLMProvider = "gemini"
    LLMProviderCloudflare  LLMProvider = "cloudflare"
    LLMProviderHuggingFace LLMProvider = "huggingface"
    LLMProviderOllama      LLMProvider = "ollama"
)
```

## API Reference

### Client

```go
type Client struct {
    // ...
}

func New(cfg *Config) (*Client, error)
func (c *Client) Summarize(ctx context.Context, text string) (string, error)
```

### Configuration

```go
type Config struct {
    Type Type
    LLM  *LLMConfig
}

type LLMConfig struct {
    Provider    LLMProvider
    Model       string
    APIKey      string
    URL         string
    BaseURL     string
    MaxTokens   int
    AccountID   string
    Credentials *Credentials
    Options     map[string]interface{}
}
```

### Summarizer Interface

```go
type Summarizer interface {
    Summarize(ctx context.Context, text string) (string, error)
}
```

## Advanced Usage

### Custom LLM Client

For maximum flexibility and testability, you can create a summarizer with a custom LLM client:

```go
import "github.com/kopexa-grc/common/llm"

// Create LLM client
llmClient, err := llm.New(llm.NewConfig(
    llm.WithOpenAI("gpt-4", "your-api-key"),
))
if err != nil {
    log.Fatal(err)
}

// Create summarizer with custom client
summarizer := summarizer.NewLLMSummarizer(llmClient)

// Use directly
summary, err := summarizer.Summarize(ctx, text)
```

### Provider-Specific Options

Each LLM provider supports specific configuration options:

#### OpenAI/Azure OpenAI
- `organization_id`: OpenAI Organization ID
- `deployment`: Azure OpenAI Deployment Name
- `api_type`: "openai", "azure", "azuread"
- `api_version`: Azure API Version
- `embedding_model`: Embedding Model Name

#### Anthropic
- `beta_header`: Beta Feature Flags
- `legacy_text_completion`: Legacy API Mode

#### Google Gemini
- `credentials`: Service Account Credentials (Path or JSON)

### Language Detection

The package automatically detects the input language and uses appropriate prompts:

- **German**: Uses German prompt for summarization
- **Other languages**: Uses English prompt as fallback

```go
// German text will be summarized in German
summary, err := client.Summarize(ctx, "Dies ist ein deutscher Text.")
```

## Error Handling

The package defines specific errors for different scenarios:

```go
var (
    ErrConfigRequired    = errors.New("config must not be nil")
    ErrLLMConfigRequired = errors.New("LLM config is required for LLM summarization")
    ErrUnsupportedType   = errors.New("unsupported summarizer type")
    ErrSentenceEmpty     = errors.New("sentence is empty after sanitization")
)
```

## Testing

### Unit Tests

```bash
go test ./summarizer/...
```

### Integration Tests

Integration tests require API keys and are skipped if not available:

```bash
# Set environment variables for integration tests
export AZURE_OPENAI_API_KEY="your-key"
export AZURE_OPENAI_ENDPOINT="your-endpoint"
export AZURE_OPENAI_DEPLOYMENT="your-deployment"

go test -v ./summarizer/...
```

### Mock Testing

The package supports easy mocking through the `LLMClient` interface:

```go
type mockLLMClient struct{}

func (m *mockLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
    return "Mock summary", nil
}

summarizer := summarizer.NewLLMSummarizer(&mockLLMClient{})
```

## Performance Considerations

### LexRank Summarization
- **Speed**: Fast, suitable for real-time applications
- **Memory**: Low memory usage
- **Quality**: Extractive, preserves original sentences

### LLM Summarization
- **Speed**: Slower, depends on API response time
- **Memory**: Higher memory usage
- **Quality**: Abstractive, generates new text

## Security

### Input Sanitization
- All input is sanitized using `bluemonday.StrictPolicy()`
- HTML tags and potentially malicious content are removed
- Empty content after sanitization results in an error

### API Key Management
- API keys should be stored securely (environment variables, secret management)
- Never hardcode API keys in source code
- Use appropriate IAM roles and permissions for cloud providers

## Examples

### Incident Report Summarization

```go
config := summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLlm),
    summarizer.WithOpenAI("gpt-4", apiKey),
)

client, err := summarizer.New(config)
if err != nil {
    log.Fatal(err)
}

incidentReport := `Critical security incident detected at 14:30 UTC. 
                   Unauthorized access attempt from IP 192.168.1.100. 
                   System administrator was notified immediately. 
                   Incident contained within 15 minutes. 
                   No data breach occurred.`

summary, err := client.Summarize(ctx, incidentReport)
```

### Technical Documentation Summarization

```go
config := summarizer.NewConfig(
    summarizer.WithType(summarizer.TypeLexrank),
)

client, err := summarizer.New(config)
if err != nil {
    log.Fatal(err)
}

docs := `This API provides endpoints for user management. 
         POST /users creates a new user. 
         GET /users/{id} retrieves user information. 
         PUT /users/{id} updates user data. 
         DELETE /users/{id} removes a user.`

summary, err := client.Summarize(ctx, docs)
```

## Integration with Other Packages

The summarizer package integrates with the `llm` package for LLM-based summarization:

```go
import (
    "github.com/kopexa-grc/common/summarizer"
    "github.com/kopexa-grc/common/llm"
)

// Use llm package directly for advanced use cases
llmClient, err := llm.New(llm.NewConfig(
    llm.WithOpenAI("gpt-4", apiKey),
))

// Use summarizer package for summarization-specific functionality
summarizer := summarizer.NewLLMSummarizer(llmClient)
```

## License

See main project license. 