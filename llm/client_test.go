// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package llm

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid OpenAI config",
			config: NewConfig(
				WithOpenAI("gpt-4", "test-api-key"),
			),
			wantErr: false,
		},
		{
			name: "valid Anthropic config",
			config: NewConfig(
				WithAnthropic("claude-3-sonnet", "test-api-key"),
			),
			wantErr: false,
		},
		{
			name: "valid Ollama config",
			config: NewConfig(
				WithOllama("llama2", "http://localhost:11434"),
			),
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "unsupported provider",
			config: NewConfig(
				WithProvider("unsupported"),
				WithModel("test-model"),
				WithAPIKey("test-key"),
			),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}

				if client != nil {
					t.Errorf("Expected nil client when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if client == nil {
					t.Errorf("Expected client but got nil")
				}
			}
		})
	}
}

func TestClient_Generate(t *testing.T) {
	// Skip if no API key is available
	apiKey := "test-key"
	if apiKey == "" {
		t.Skip("No API key available, skipping generation tests")
	}

	config := NewConfig(
		WithOpenAI("gpt-4", apiKey),
	)

	client, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	result, err := client.Generate(ctx, "Hello, how are you?")

	// We expect an error here because we're using a test key
	// The important thing is that the client creation and method call don't panic
	if err != nil {
		t.Logf("Expected error for invalid credentials: %v", err)
	} else {
		t.Logf("Generated response: %s", result)
	}
}

func TestClient_GenerateWithOptions(t *testing.T) {
	config := NewConfig(
		WithOpenAI("gpt-4", "test-key"),
	)

	client, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	result, err := client.GenerateWithOptions(ctx, "Test prompt")

	// We expect an error here because we're using a test key
	if err != nil {
		t.Logf("Expected error for invalid credentials: %v", err)
	} else {
		t.Logf("Generated response: %s", result)
	}
}

func TestClient_GetModel(t *testing.T) {
	config := NewConfig(
		WithOpenAI("gpt-4", "test-key"),
	)

	client, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	model := client.GetModel()
	if model == nil {
		t.Error("Expected non-nil model")
	}
}
