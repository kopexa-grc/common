// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import (
	"context"
	"errors"
	"strings"

	"github.com/didasy/tldr"
)

// ErrInvalidMaxSentences is returned when maxSentences is less than 1
var ErrInvalidMaxSentences = errors.New("maxSentences must be at least 1")

// lexRankSummarizer implements the LexRank algorithm for extractive summarization
type lexRankSummarizer struct {
	maxSentences int
}

// newLexRankSummarizer creates a new LexRank summarizer with the specified configuration
func newLexRankSummarizer(maxSentences int) (*lexRankSummarizer, error) {
	if maxSentences < 1 {
		return nil, ErrInvalidMaxSentences
	}

	return &lexRankSummarizer{
		maxSentences: maxSentences,
	}, nil
}

// Summarize performs extractive summarization using the LexRank algorithm
//
// The LexRank algorithm ranks sentences based on their centrality in the document graph.
// It identifies the most important sentences by analyzing the similarity between sentences
// and their connections in the document.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - text: The input text to summarize
//
// Returns:
//   - The summarized text as a string
//   - An error if summarization fails
func (l *lexRankSummarizer) Summarize(ctx context.Context, text string) (string, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Validate input
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "" {
		return "", ErrSentenceEmpty
	}

	// Create LexRank summarizer instance
	summarizer := tldr.New()

	// Perform summarization
	sentences, err := summarizer.Summarize(trimmedText, l.maxSentences)
	if err != nil {
		return "", err
	}

	// Join sentences into final summary
	summary := strings.Join(sentences, " ")

	// Handle edge case: if the algorithm returns an empty summary
	// (can happen with very short texts), return the original text
	if len(summary) == 0 {
		summary = trimmedText
	}

	return summary, nil
}

// GetMaxSentences returns the maximum number of sentences for summarization
func (l *lexRankSummarizer) GetMaxSentences() int {
	return l.maxSentences
}
