// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import (
	"context"
	"fmt"

	"github.com/abadojack/whatlanggo"

	"github.com/kopexa-grc/common/llm"
)

const (
	promptEN = `Summarize the following text in English. Be brief, concise and precise.\n\n%s\n`
	promptDE = `Fasse den folgenden Text auf Deutsch zusammen. Sei kurz, prägnant und präzise.\n\n%s\n`
)

// LLMClient interface for LLM-based text generation
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// LLMSummarizer implements summarization using LLM clients
type LLMSummarizer struct {
	llmClient LLMClient
}

// NewLLMSummarizer creates a summarizer from an existing LLMClient
func NewLLMSummarizer(client LLMClient) *LLMSummarizer {
	return &LLMSummarizer{llmClient: client}
}

// NewLLMSummarizerFromConfig is a convenience constructor that builds the client from a summarizer Config
func NewLLMSummarizerFromConfig(cfg Config) (*LLMSummarizer, error) {
	if cfg.LLM == nil {
		return nil, ErrLLMConfigRequired
	}

	var llmConfig *llm.Config

	switch cfg.LLM.Provider {
	case LLMProviderOpenAI:
		llmConfig = llm.NewConfig(
			llm.WithOpenAI(cfg.LLM.Model, cfg.LLM.APIKey,
				llm.WithURL(cfg.LLM.URL),
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
				llm.WithOption("organization_id", cfg.LLM.Options["organization_id"]),
				llm.WithOption("deployment", cfg.LLM.Options["deployment"]),
				llm.WithOption("api_type", cfg.LLM.Options["api_type"]),
				llm.WithOption("api_version", cfg.LLM.Options["api_version"]),
				llm.WithOption("embedding_model", cfg.LLM.Options["embedding_model"]),
			),
		)
	case LLMProviderAnthropic:
		llmConfig = llm.NewConfig(
			llm.WithAnthropic(cfg.LLM.Model, cfg.LLM.APIKey,
				llm.WithBaseURL(cfg.LLM.BaseURL),
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
				llm.WithOption("beta_header", cfg.LLM.Options["beta_header"]),
				llm.WithOption("legacy_text_completion", cfg.LLM.Options["legacy_text_completion"]),
			),
		)
	case LLMProviderMistral:
		llmConfig = llm.NewConfig(
			llm.WithMistral(cfg.LLM.Model, cfg.LLM.APIKey, cfg.LLM.URL,
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
			),
		)
	case LLMProviderGemini:
		llmConfig = llm.NewConfig(
			llm.WithGemini(cfg.LLM.Model,
				llm.WithAPIKey(cfg.LLM.APIKey),
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
				llm.WithCredentials(cfg.LLM.Credentials.Path, cfg.LLM.Credentials.JSON),
			),
		)
	case LLMProviderHuggingFace:
		llmConfig = llm.NewConfig(
			llm.WithHuggingFace(cfg.LLM.Model, cfg.LLM.APIKey, cfg.LLM.URL,
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
			),
		)
	case LLMProviderOllama:
		llmConfig = llm.NewConfig(
			llm.WithOllama(cfg.LLM.Model, cfg.LLM.URL,
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
			),
		)
	case LLMProviderCloudflare:
		llmConfig = llm.NewConfig(
			llm.WithCloudflare(cfg.LLM.Model, cfg.LLM.APIKey, cfg.LLM.AccountID,
				llm.WithURL(cfg.LLM.URL),
				llm.WithMaxTokens(cfg.LLM.MaxTokens),
			),
		)
	default:
		return nil, fmt.Errorf("%w: %s", llm.ErrUnsupportedProvider, cfg.LLM.Provider)
	}

	client, err := llm.New(llmConfig)
	if err != nil {
		return nil, err
	}

	return NewLLMSummarizer(client), nil
}

// Summarize returns a shortened version of the provided string using the selected llm
func (l *LLMSummarizer) Summarize(ctx context.Context, s string) (string, error) {
	langInfo := whatlanggo.Detect(s)

	var prompt string

	switch langInfo.Lang.String() {
	case "German":
		prompt = promptDE
	default:
		prompt = promptEN
	}

	return l.llmClient.Generate(ctx, fmt.Sprintf(prompt, s))
}
