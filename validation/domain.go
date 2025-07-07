// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package validation provides domain validation utilities for URL and network operations.
//
// This package implements comprehensive validation functions following the Google API
// Design Guide principles for input validation, error handling, and network operations.
// The package provides both syntactic validation (URL format) and semantic validation
// (network reachability) with proper error categorization and user-friendly messages.
//
// The validation functions are designed to be:
//   - Fast and efficient for high-throughput applications
//   - Secure against common attack vectors (URL injection, DNS attacks)
//   - Configurable for different deployment environments
//   - Well-documented with clear error messages
//
// Example usage:
//
//	// Basic URL validation
//	if err := validation.IsValidURL("https://example.com"); err != nil {
//		log.Printf("Invalid URL: %v", err)
//	}
//
//	// URL reachability check with timeout
//	if err := validation.CheckURLReachability("https://api.example.com"); err != nil {
//		log.Printf("URL not reachable: %v", err)
//	}
//
// The package supports both HTTP and HTTPS schemes and includes protection against
// common security issues such as overly long URLs, invalid domain names, and
// network timeouts.
package validation

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/kopexa-grc/common/errors"
)

// Error codes for domain validation operations.
//
// These error codes follow the Google API Design Guide pattern of using
// descriptive, hierarchical error codes that can be easily categorized
// and handled by clients.
const (
	// ErrCodeInvalidURL indicates that the provided URL is syntactically invalid.
	// This includes malformed URLs, unsupported schemes, or invalid domain names.
	ErrCodeInvalidURL = "VALIDATION_INVALID_URL"

	// ErrCodeURLTooLong indicates that the provided URL exceeds the maximum allowed length.
	// This prevents potential buffer overflow attacks and ensures reasonable URL sizes.
	ErrCodeURLTooLong = "VALIDATION_URL_TOO_LONG"

	// ErrCodeEmptyURL indicates that an empty URL was provided.
	// This prevents processing of empty strings that could cause unexpected behavior.
	ErrCodeEmptyURL = "VALIDATION_EMPTY_URL"

	// ErrCodeUnsupportedScheme indicates that the URL uses an unsupported protocol scheme.
	// Currently only HTTP and HTTPS are supported for security reasons.
	ErrCodeUnsupportedScheme = "VALIDATION_UNSUPPORTED_SCHEME"

	// ErrCodeInvalidDomain indicates that the URL contains an invalid domain name.
	// This includes domains with invalid characters or improper formatting.
	ErrCodeInvalidDomain = "VALIDATION_INVALID_DOMAIN"

	// ErrCodeHostNotFound indicates that the domain name could not be resolved.
	// This includes DNS resolution failures and non-existent domains.
	ErrCodeHostNotFound = "VALIDATION_HOST_NOT_FOUND"

	// ErrCodeRequestCreationFailed indicates that an HTTP request could not be created.
	// This is typically due to malformed URLs or invalid request parameters.
	ErrCodeRequestCreationFailed = "VALIDATION_REQUEST_CREATION_FAILED"

	// ErrCodeHTTPRequestFailed indicates that the HTTP request failed to complete.
	// This includes network timeouts, connection failures, and transport errors.
	ErrCodeHTTPRequestFailed = "VALIDATION_HTTP_REQUEST_FAILED"

	// ErrCodeNonSuccessStatusCode indicates that the HTTP request completed but
	// returned a non-success status code (4xx or 5xx).
	ErrCodeNonSuccessStatusCode = "VALIDATION_NON_SUCCESS_STATUS_CODE"
)

// Configuration constants for URL validation.
//
// These constants define the operational parameters for URL validation
// and can be adjusted based on deployment requirements and security policies.
const (
	// MaxURLLength defines the maximum allowed length for URLs.
	// This prevents potential buffer overflow attacks and ensures reasonable
	// URL sizes. The value of 2048 characters follows common web standards
	// and browser limitations.
	MaxURLLength = 2048

	// DefaultHTTPTimeout defines the default timeout for HTTP operations.
	// This timeout applies to both DNS resolution and HTTP request completion.
	// The value of 5 seconds provides a good balance between responsiveness
	// and reliability for most network conditions.
	DefaultHTTPTimeout = 5 * time.Second

	// DefaultUserAgent defines the User-Agent string used for HTTP requests.
	// This helps identify the application making requests and can be useful
	// for debugging and monitoring purposes.
	DefaultUserAgent = "Kopexa-Validation/1.0"

	DialTimeout           = 30 * time.Second
	DialKeepAlive         = 30 * time.Second
	TLSHandshakeTimeout   = 10 * time.Second
	ResponseHeaderTimeout = 10 * time.Second
	IdleConnTimeout       = 90 * time.Second
)

// Supported URL schemes for validation.
//
// Only HTTP and HTTPS schemes are supported for security reasons.
// Additional schemes can be added here if needed for specific use cases.
var supportedSchemes = []string{"http", "https"}

// Domain name validation regex pattern.
//
// This regex pattern validates domain names according to RFC 1035 standards
// with some additional restrictions for security and practicality:
//   - Case-insensitive matching
//   - Allows letters, numbers, and hyphens
//   - Requires at least one dot separator
//   - Prevents leading/trailing hyphens in labels
//   - Supports optional trailing dot for root domains
//   - Prevents single-label domains (must have at least one dot)
var domainRegexp = regexp.MustCompile(`^(?i)[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+\.?$`)

// IsValidURL validates the syntax and format of a URL string.
//
// This function performs comprehensive URL validation including:
//   - Length validation to prevent buffer overflow attacks
//   - Scheme validation to ensure only supported protocols are used
//   - Domain name validation using RFC-compliant regex patterns
//   - Basic URL structure validation
//
// The function supports URLs with or without scheme prefixes. If no scheme
// is provided, "http://" is assumed as the default.
//
// Returns nil if the URL is valid, or an error with appropriate error code
// and descriptive message if validation fails.
//
// Example:
//
//	err := IsValidURL("https://example.com/path")
//	if err != nil {
//		// Handle validation error
//	}
//
//	err = IsValidURL("example.com") // Assumes http://
//	if err != nil {
//		// Handle validation error
//	}
func IsValidURL(inputURL string) error {
	// Validate URL length to prevent potential attacks
	if inputURL == "" {
		return errors.New(ErrCodeEmptyURL, "URL cannot be empty")
	}

	if len(inputURL) > MaxURLLength {
		return errors.New(ErrCodeURLTooLong, fmt.Sprintf("URL length %d exceeds maximum allowed length of %d", len(inputURL), MaxURLLength))
	}

	// Perform detailed URL validation
	if err := validateURLSyntax(inputURL); err != nil {
		return err
	}

	return nil
}

// validateURLSyntax performs detailed URL syntax validation.
//
// This internal function handles the core URL parsing and validation logic,
// including scheme validation, domain name validation, and structural checks.
// It is separated from the public interface to allow for better testing
// and code organization.
func validateURLSyntax(inputURL string) error {
	// Parse the URL to validate its structure
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return errors.New(ErrCodeInvalidURL, fmt.Sprintf("URL parsing failed: %v", err))
	}

	// Handle URLs without scheme by adding default HTTP scheme
	if parsedURL.Scheme == "" {
		parsedURL, err = url.Parse("http://" + inputURL)
		if err != nil {
			return errors.New(ErrCodeInvalidURL, fmt.Sprintf("URL parsing with default scheme failed: %v", err))
		}
	}

	// Validate that the host is present
	if parsedURL.Host == "" {
		return errors.New(ErrCodeInvalidURL, "URL must contain a valid host")
	}

	// Validate the URL scheme
	if !slices.Contains(supportedSchemes, parsedURL.Scheme) {
		return errors.New(ErrCodeUnsupportedScheme, fmt.Sprintf("Unsupported URL scheme '%s'. Only %v are supported", parsedURL.Scheme, supportedSchemes))
	}

	// Validate the domain name format
	if !isValidDomain(parsedURL.Host) {
		return errors.New(ErrCodeInvalidDomain, fmt.Sprintf("Invalid domain name '%s'", parsedURL.Host))
	}

	return nil
}

// isValidDomain validates a domain name using regex pattern matching.
//
// This function checks that the domain name follows RFC 1035 standards
// and includes additional security restrictions to prevent common
// attack vectors such as domain name spoofing.
func isValidDomain(host string) bool {
	// Remove port if present for domain validation
	hostname := host
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		hostname = host[:colonIndex]
	}

	return domainRegexp.MatchString(hostname)
}

// CheckURLReachability validates that a URL is both syntactically valid and
// network reachable.
//
// This function performs a comprehensive validation that includes:
//   - URL syntax validation using IsValidURL
//   - DNS resolution to verify the domain exists
//   - HTTP HEAD request to verify the endpoint is accessible
//   - Status code validation to ensure the service is operational
//
// The function uses configurable timeouts to prevent hanging operations
// and includes proper error categorization for different failure modes.
//
// Network operations are performed with appropriate timeouts and user agent
// identification to ensure reliable and traceable requests.
//
// Returns nil if the URL is reachable, or an error with appropriate error
// code and descriptive message if the URL cannot be reached.
//
// Example:
//
//	err := CheckURLReachability("https://api.example.com")
//	if err != nil {
//		// Handle reachability error
//		switch {
//		case errors.Is(err, ErrCodeHostNotFound):
//			// DNS resolution failed
//		case errors.Is(err, ErrCodeHTTPRequestFailed):
//			// Network request failed
//		case errors.Is(err, ErrCodeNonSuccessStatusCode):
//			// Service returned error status
//		}
//	}
func CheckURLReachability(rawURL string) error {
	// First validate the URL syntax
	if err := IsValidURL(rawURL); err != nil {
		return err
	}

	// Parse the URL for network operations
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.New(ErrCodeInvalidURL, fmt.Sprintf("URL parsing failed during reachability check: %v", err))
	}

	// Perform DNS resolution to verify the domain exists
	if err := validateDNSResolution(parsedURL.Hostname()); err != nil {
		return err
	}

	// Perform HTTP reachability check
	if err := validateHTTPReachability(rawURL); err != nil {
		return err
	}

	return nil
}

// validateDNSResolution performs DNS resolution for a hostname.
//
// This function verifies that the domain name can be resolved to an IP address,
// which is a prerequisite for any network communication. DNS resolution failures
// typically indicate either network connectivity issues or non-existent domains.
func validateDNSResolution(hostname string) error {
	// Perform DNS lookup with timeout
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHTTPTimeout)
	defer cancel()

	// Use a custom resolver with timeout
	resolver := &net.Resolver{
		PreferGo: true,
	}

	_, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		return errors.New(ErrCodeHostNotFound, fmt.Sprintf("DNS resolution failed for '%s': %v", hostname, err))
	}

	return nil
}

// validateHTTPReachability performs an HTTP HEAD request to verify endpoint accessibility.
//
// This function sends an HTTP HEAD request to the URL to verify that:
//   - The endpoint is accessible over the network
//   - The service is responding to HTTP requests
//   - The service returns a success status code
//
// HEAD requests are used instead of GET requests to minimize bandwidth usage
// while still verifying that the service is operational.
func validateHTTPReachability(rawURL string) error {
	// Create context with timeout for the HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHTTPTimeout)
	defer cancel()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, http.NoBody)
	if err != nil {
		return errors.New(ErrCodeRequestCreationFailed, fmt.Sprintf("Failed to create HTTP request: %v", err))
	}

	// Set user agent for request identification
	req.Header.Set("User-Agent", DefaultUserAgent)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: DefaultHTTPTimeout,
		Transport: &http.Transport{
			// Disable keep-alive to ensure fresh connections
			DisableKeepAlives: true,
			// Set reasonable timeouts for connection establishment
			DialContext: (&net.Dialer{
				Timeout:   DialTimeout,
				KeepAlive: DialKeepAlive,
			}).DialContext,
			// Set timeouts for TLS handshake
			TLSHandshakeTimeout: TLSHandshakeTimeout,
			// Set timeouts for response headers
			ResponseHeaderTimeout: ResponseHeaderTimeout,
			// Set timeouts for idle connections
			IdleConnTimeout: IdleConnTimeout,
		},
	}

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(ErrCodeHTTPRequestFailed, fmt.Sprintf("HTTP request failed: %v", err))
	}
	defer resp.Body.Close()

	// Validate the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return errors.New(ErrCodeNonSuccessStatusCode, fmt.Sprintf("HTTP request returned non-success status code: %d", resp.StatusCode))
	}

	return nil
}
