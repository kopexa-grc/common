// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package validation

import (
	"testing"
	"time"

	"github.com/kopexa-grc/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorCode   string
		description string
	}{
		{
			name:        "valid HTTPS URL",
			input:       "https://example.com",
			expectError: false,
			description: "should accept valid HTTPS URLs",
		},
		{
			name:        "valid HTTP URL",
			input:       "http://example.com",
			expectError: false,
			description: "should accept valid HTTP URLs",
		},
		{
			name:        "valid URL without scheme",
			input:       "example.com",
			expectError: false,
			description: "should accept URLs without scheme (assumes HTTP)",
		},
		{
			name:        "valid URL with path",
			input:       "https://example.com/path",
			expectError: false,
			description: "should accept URLs with paths",
		},
		{
			name:        "valid URL with query parameters",
			input:       "https://example.com?param=value",
			expectError: false,
			description: "should accept URLs with query parameters",
		},
		{
			name:        "valid URL with port",
			input:       "https://example.com:8080",
			expectError: false,
			description: "should accept URLs with ports",
		},
		{
			name:        "valid subdomain",
			input:       "https://sub.example.com",
			expectError: false,
			description: "should accept URLs with subdomains",
		},
		{
			name:        "valid multi-level domain",
			input:       "https://sub.sub2.example.com",
			expectError: false,
			description: "should accept URLs with multiple subdomain levels",
		},
		{
			name:        "empty URL",
			input:       "",
			expectError: true,
			errorCode:   ErrCodeEmptyURL,
			description: "should reject empty URLs",
		},
		{
			name:        "URL too long",
			input:       "https://" + string(make([]byte, MaxURLLength+1)),
			expectError: true,
			errorCode:   ErrCodeURLTooLong,
			description: "should reject URLs exceeding maximum length",
		},
		{
			name:        "unsupported scheme FTP",
			input:       "ftp://example.com",
			expectError: true,
			errorCode:   ErrCodeUnsupportedScheme,
			description: "should reject unsupported URL schemes",
		},
		{
			name:        "unsupported scheme SSH",
			input:       "ssh://example.com",
			expectError: true,
			errorCode:   ErrCodeUnsupportedScheme,
			description: "should reject unsupported URL schemes",
		},
		{
			name:        "invalid domain with spaces",
			input:       "https://example .com",
			expectError: true,
			errorCode:   ErrCodeInvalidURL,
			description: "should reject domains with invalid characters",
		},
		{
			name:        "invalid domain with special characters",
			input:       "https://example@.com",
			expectError: true,
			errorCode:   ErrCodeInvalidDomain,
			description: "should reject domains with special characters",
		},
		{
			name:        "domain starting with hyphen",
			input:       "https://-example.com",
			expectError: true,
			errorCode:   ErrCodeInvalidDomain,
			description: "should reject domains starting with hyphen",
		},
		{
			name:        "domain ending with hyphen",
			input:       "https://example-.com",
			expectError: true,
			errorCode:   ErrCodeInvalidDomain,
			description: "should reject domains ending with hyphen",
		},
		{
			name:        "single label domain",
			input:       "https://example",
			expectError: true,
			errorCode:   ErrCodeInvalidDomain,
			description: "should reject single label domains",
		},
		{
			name:        "malformed URL",
			input:       "not-a-url",
			expectError: true,
			errorCode:   ErrCodeInvalidDomain,
			description: "should reject malformed URLs",
		},
		{
			name:        "URL with invalid host",
			input:       "https://",
			expectError: true,
			errorCode:   ErrCodeInvalidURL,
			description: "should reject URLs without host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidURL(tt.input)

			if tt.expectError {
				require.Error(t, err, tt.description)

				if tt.errorCode != "" {
					assert.Equal(t, tt.errorCode, string(errors.Code(err)), "error code should match expected")
				}
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

func TestIsValidURL_EdgeCases(t *testing.T) {
	t.Run("URL with trailing dot", func(t *testing.T) {
		// RFC 1035 allows trailing dots for root domains
		err := IsValidURL("https://example.com.")
		assert.NoError(t, err, "should accept URLs with trailing dot")
	})

	t.Run("URL with multiple dots", func(t *testing.T) {
		err := IsValidURL("https://example..com")
		assert.Error(t, err, "should reject URLs with consecutive dots")
		assert.Equal(t, ErrCodeInvalidDomain, string(errors.Code(err)))
	})

	t.Run("URL with underscore in domain", func(t *testing.T) {
		err := IsValidURL("https://example_test.com")
		assert.Error(t, err, "should reject domains with underscores")
		assert.Equal(t, ErrCodeInvalidDomain, string(errors.Code(err)))
	})

	t.Run("URL with uppercase letters", func(t *testing.T) {
		err := IsValidURL("https://EXAMPLE.COM")
		assert.NoError(t, err, "should accept domains with uppercase letters (case insensitive)")
	})
}

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"valid domain", "example.com", true},
		{"valid subdomain", "sub.example.com", true},
		{"valid multi-level", "sub1.sub2.example.com", true},
		{"valid with numbers", "example123.com", true},
		{"valid with hyphens", "example-domain.com", true},
		{"valid trailing dot", "example.com.", true},
		{"valid with port", "example.com:8080", true},
		{"invalid single label", "example", false},
		{"invalid with spaces", "example .com", false},
		{"invalid with special chars", "example@.com", false},
		{"invalid starting hyphen", "-example.com", false},
		{"invalid ending hyphen", "example-.com", false},
		{"invalid consecutive dots", "example..com", false},
		{"invalid with underscore", "example_test.com", false},
		{"empty host", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidDomain(tt.host)
			assert.Equal(t, tt.expected, result, "domain validation should match expected result")
		})
	}
}

func TestValidateURLSyntax(t *testing.T) {
	t.Run("valid URL parsing", func(t *testing.T) {
		err := validateURLSyntax("https://example.com")
		assert.NoError(t, err, "should parse valid URLs without error")
	})

	t.Run("URL parsing failure", func(t *testing.T) {
		// Create a URL that will fail parsing
		invalidURL := string([]byte{0x00, 0x01, 0x02}) // Invalid UTF-8 sequence
		err := validateURLSyntax(invalidURL)
		assert.Error(t, err, "should fail to parse invalid URLs")
		assert.Equal(t, ErrCodeInvalidURL, string(errors.Code(err)))
	})

	t.Run("URL with default scheme", func(t *testing.T) {
		err := validateURLSyntax("example.com")
		assert.NoError(t, err, "should add default HTTP scheme to URLs without scheme")
	})

	t.Run("URL with default scheme parsing failure", func(t *testing.T) {
		// Create a URL that will fail parsing even with default scheme
		invalidURL := "http://" + string([]byte{0x00, 0x01, 0x02})
		err := validateURLSyntax(invalidURL)
		assert.Error(t, err, "should fail to parse URLs with default scheme")
		assert.Equal(t, ErrCodeInvalidURL, string(errors.Code(err)))
	})
}

func TestCheckURLReachability(t *testing.T) {
	// Note: These tests require network connectivity and may be flaky
	// In a real enterprise environment, these would be mocked or run in integration tests
	t.Run("valid reachable URL", func(t *testing.T) {
		// Use a well-known, reliable service for testing
		err := CheckURLReachability("https://httpbin.org/status/200")
		// This test may fail if network is unavailable, so we don't assert on the result
		// In a real test environment, this would be mocked
		t.Logf("Reachability check result: %v", err)
	})

	t.Run("invalid URL", func(t *testing.T) {
		err := CheckURLReachability("not-a-url")
		assert.Error(t, err, "should fail for invalid URLs")
		// Der Fehler kann je nach Parsing INVALID_DOMAIN, EMPTY_URL oder HOST_NOT_FOUND sein
		errorCode := string(errors.Code(err))
		assert.Contains(t, []string{ErrCodeEmptyURL, ErrCodeHostNotFound, ErrCodeInvalidDomain}, errorCode,
			"error code should be empty URL, host not found, or invalid domain")
	})

	t.Run("empty URL", func(t *testing.T) {
		err := CheckURLReachability("")
		assert.Error(t, err, "should fail for empty URLs")
		assert.Equal(t, ErrCodeEmptyURL, string(errors.Code(err)))
	})

	t.Run("unsupported scheme", func(t *testing.T) {
		err := CheckURLReachability("ftp://example.com")
		assert.Error(t, err, "should fail for unsupported schemes")
		assert.Equal(t, ErrCodeUnsupportedScheme, string(errors.Code(err)))
	})
}

func TestValidateDNSResolution(t *testing.T) {
	t.Run("valid hostname", func(t *testing.T) {
		// Use a well-known domain for testing
		err := validateDNSResolution("google.com")
		// This test may fail if DNS is unavailable, so we don't assert on the result
		t.Logf("DNS resolution result: %v", err)
	})

	t.Run("invalid hostname", func(t *testing.T) {
		// Use a domain that should not exist
		err := validateDNSResolution("this-domain-should-not-exist-12345.com")
		// This test may pass if the domain is registered, so we don't assert on the result
		t.Logf("DNS resolution result for non-existent domain: %v", err)
	})
}

func TestValidateHTTPReachability(t *testing.T) {
	t.Run("valid HTTP endpoint", func(t *testing.T) {
		// Use a reliable test service
		err := validateHTTPReachability("https://httpbin.org/status/200")
		// This test may fail if network is unavailable, so we don't assert on the result
		t.Logf("HTTP reachability result: %v", err)
	})

	t.Run("invalid URL format", func(t *testing.T) {
		err := validateHTTPReachability("not-a-url")
		assert.Error(t, err, "should fail for invalid URL format")
		// The error could be either request creation failed or HTTP request failed
		errorCode := string(errors.Code(err))
		assert.True(t, errorCode == ErrCodeRequestCreationFailed || errorCode == ErrCodeHTTPRequestFailed,
			"error code should be either request creation failed or HTTP request failed")
	})
}

func TestConstants(t *testing.T) {
	t.Run("MaxURLLength", func(t *testing.T) {
		assert.Greater(t, MaxURLLength, 0, "MaxURLLength should be positive")
		assert.LessOrEqual(t, MaxURLLength, 8192, "MaxURLLength should be reasonable")
	})

	t.Run("DefaultHTTPTimeout", func(t *testing.T) {
		assert.Greater(t, DefaultHTTPTimeout, time.Duration(0), "DefaultHTTPTimeout should be positive")
		assert.LessOrEqual(t, DefaultHTTPTimeout, 30*time.Second, "DefaultHTTPTimeout should be reasonable")
	})

	t.Run("DefaultUserAgent", func(t *testing.T) {
		assert.NotEmpty(t, DefaultUserAgent, "DefaultUserAgent should not be empty")
		assert.Contains(t, DefaultUserAgent, "Kopexa", "DefaultUserAgent should identify the application")
	})

	t.Run("supported schemes", func(t *testing.T) {
		assert.Contains(t, supportedSchemes, "http", "should support HTTP scheme")
		assert.Contains(t, supportedSchemes, "https", "should support HTTPS scheme")
		assert.Len(t, supportedSchemes, 2, "should only support HTTP and HTTPS")
	})
}

func TestErrorCodes(t *testing.T) {
	// Test that all error codes are properly defined
	errorCodes := []string{
		ErrCodeInvalidURL,
		ErrCodeURLTooLong,
		ErrCodeEmptyURL,
		ErrCodeUnsupportedScheme,
		ErrCodeInvalidDomain,
		ErrCodeHostNotFound,
		ErrCodeRequestCreationFailed,
		ErrCodeHTTPRequestFailed,
		ErrCodeNonSuccessStatusCode,
	}

	for _, code := range errorCodes {
		t.Run(code, func(t *testing.T) {
			assert.NotEmpty(t, code, "error code should not be empty")
			assert.Contains(t, code, "VALIDATION_", "error code should have VALIDATION_ prefix")
		})
	}
}
