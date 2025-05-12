# Blob Storage Package

## Motivation

The Blob Storage Package provides a unified interface for working with various blob storage services. It is based on the [gocloud.dev/blob](https://gocloud.dev/howto/blob/) package and extends it with Kopexa-specific functionality.

## Why this Package?

1. **Unified Interface**: The package abstracts the complexity of different blob storage services (such as Azure Blob Storage) behind a simple, unified API.

2. **Security**: Implements secure access controls and supports signed URLs for temporary access.

3. **Flexibility**: Supports different access levels (public/private) and container types.

4. **Error Handling**: Integrates with the Kopexa Error Handling System for consistent error handling.

## Features

- Simple bucket operations (Read, Write, Delete)
- Support for signed URLs
- Copy operations between blobs
- Thread-safe implementation
- UTF-8 validation for keys
- Azure Blob Storage integration

## Usage

```go
config := &blob.Config{
    Azure: blob.AzureConfig{
        AccountName: "your-account",
        AccountKey:  "your-key",
        Endpoint:    "your-endpoint",
    },
}

provider, err := blob.New(config)
if err != nil {
    // Handle error
}

// Open a public bucket
publicBucket, err := provider.Public()
if err != nil {
    // Handle error
}

// Open a private bucket for a space
spaceBucket, err := provider.Space("space-id")
if err != nil {
    // Handle error
}
``` 