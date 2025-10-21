This code was stolen from: https://github.com/theopenlane/iam/tree/main/fgax

# FGA (Fine-Grained Authorization) Package

The FGA package provides a type-safe and idiomatic Go implementation for interacting with OpenFGA (Fine-Grained Authorization). It implements a fluent interface for managing permissions in a type-safe manner.

## Features

- **Type Safety**: All operations are strongly typed to prevent runtime errors
- **Fluent Interface**: Uses the builder pattern for clear and readable API calls
- **Efficiency**: Minimizes allocations and network calls
- **Reliability**: Handles errors gracefully and provides detailed error information
- **Configurability**: Supports various options like StoreID and duplicate handling

## Installation

```bash
go get github.com/kopexa-grc/common/fga
```

## Usage

```go
// Create client
client, err := fga.NewClient("https://api.openfga.example",
    fga.WithStoreID("store123"),
    fga.WithIgnoreDuplicateKeyError(true),
)
if err != nil {
    log.Fatal(err)
}

// Check permission
hasAccess, err := client.Has().
    User("user123").
    Capability("viewer").
    In("document", "doc123").
    Check(ctx)

// Grant permission
err = client.Grant().
    User("user123").
    Relation("viewer").
    To("document", "doc123").
    Apply(ctx)
```

## API Overview

### Client

- `NewClient(host string, opts ...Option)`: Creates a new FGA client
- `Has()`: Starts a permission check
- `Grant()`: Starts a permission grant
- `Revoke()`: Starts a permission revocation
- `ListTuples()`: Lists tuples based on filters
- `WriteTupleKeys()`: Writes or deletes multiple tuples

### Options

- `WithStoreID(storeID string)`: Sets the store ID
- `WithIgnoreDuplicateKeyError(ignore bool)`: Configures duplicate error handling
