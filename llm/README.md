# LLM Package

The `llm` package provides a unified interface for various Large Language Model (LLM) providers. It abstracts the complexity of different APIs and offers consistent configuration through the Options Pattern.

## Features

- **Multi-Provider Support**: OpenAI, Anthropic, Mistral, Google Gemini, HuggingFace, Ollama, Cloudflare
- **Unified API**: Consistent interface for all providers
- **Flexible Configuration**: Options Pattern for type-safe configuration
- **Azure OpenAI Support**: Complete support for Azure OpenAI
- **Extensible**: Easy integration of new providers

## Installation

```bash
go get github.com/kopexa-grc/common/llm
```

## Quick Start

### Simple Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kopexa-grc/common/llm"
)

func main() {
    // Create OpenAI client
    client, err := llm.New(llm.NewConfig(
        llm.WithOpenAI("gpt-4", "your-api-key"),
    ))
    if err != nil {
        log.Fatal(err)
    }

    // Generate text
    ctx := context.Background()
    result, err := client.Generate(ctx, "Explain Go to me in one sentence.")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result)
}
```

### Azure OpenAI

```go
client, err := llm.New(llm.NewConfig(
    llm.WithOpenAI("gpt-4", "your-api-key",
        llm.WithURL("https://your-resource.openai.azure.com/"),
        llm.WithOption("deployment", "your-deployment-name"),
        llm.WithOption("api_type", "azure"),
        llm.WithOption("api_version", "2023-05-15"),
    ),
))
```

### Anthropic Claude

```go
client, err := llm.New(llm.NewConfig(
    llm.WithAnthropic("claude-3-sonnet", "your-api-key",
        llm.WithMaxTokens(1000),
    ),
))
```

### Local Ollama

```go
client, err := llm.New(llm.NewConfig(
    llm.WithOllama("llama2", "http://localhost:11434"),
))
```

## Configuration

### Options Pattern

The package uses the Options Pattern for flexible and type-safe configuration:

```go
config := llm.NewConfig(
    llm.WithProvider(llm.ProviderOpenAI),
    llm.WithModel("gpt-4"),
    llm.WithAPIKey("your-api-key"),
    llm.WithMaxTokens(1000),
    llm.WithOption("temperature", 0.7),
)
```

### Provider-specific Options

Each provider supports special options:

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

## API Reference

### Client

```go
type Client struct {
    // ...
}

func New(cfg *Config) (*Client, error)
func (c *Client) Generate(ctx context.Context, prompt string) (string, error)
func (c *Client) GenerateWithOptions(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
func (c *Client) GetModel() llms.Model
```

### Configuration

```go
type Config struct {
    Provider    Provider
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

## Error Handling

The package defines specific errors:

```go
var (
    ErrConfigRequired      = errors.New("config must not be nil")
    ErrUnsupportedProvider = errors.New("unsupported llm provider")
    ErrInvalidCredentials  = errors.New("invalid credentials provided")
)
```

## Testing

```bash
go test ./llm/...
```

## Integration with Other Packages

The `llm` package is used by other packages like `summarizer`:

```go
import "github.com/kopexa-grc/common/llm"

// LLM client for summarization
client, err := llm.New(llm.NewConfig(
    llm.WithOpenAI("gpt-4", apiKey),
))
```

## License

See main project license. 