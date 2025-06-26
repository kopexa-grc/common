// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
)

// DefaultLexRankSentences is the default number of sentences for LexRank summarization
const DefaultLexRankSentences = 3

// summarizer is the interface for all summarizer implementations
// (LexRank, LLM, ...)
type summarizer interface {
	Summarize(context.Context, string) (string, error)
}

// Client is the main entry point for summarization
// It selects the correct summarizer based on the config
// and sanitizes input/output.
type Client struct {
	impl      summarizer
	sanitizer *bluemonday.Policy
}

// New creates a new summarizer client with the given config
func New(cfg *Config) (*Client, error) {
	sanitizer := bluemonday.StrictPolicy()

	if cfg == nil {
		return nil, ErrConfigRequired
	}

	var impl summarizer

	var err error

	switch cfg.Type {
	case TypeLexrank:
		// Default: 3 Sätze, kann später erweitert werden
		impl, err = newLexRankSummarizer(DefaultLexRankSentences)
		if err != nil {
			return nil, err
		}
	case TypeLlm:
		impl, err = NewLLMSummarizerFromConfig(*cfg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedType
	}

	return &Client{
		impl:      impl,
		sanitizer: sanitizer,
	}, nil
}

// Summarize cleans the input, runs the summarizer, and sanitizes the output
func (s *Client) Summarize(ctx context.Context, sentence string) (string, error) {
	cleanInput := s.sanitizer.Sanitize(sentence)
	if cleanInput == "" {
		return "", ErrSentenceEmpty
	}

	summary, err := s.impl.Summarize(ctx, cleanInput)
	if err != nil {
		return "", err
	}

	return s.sanitizer.Sanitize(summary), nil
}
