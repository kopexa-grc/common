# FGA (Fine-Grained Authorization)

This package provides a thin wrapper around **[OpenFGA](https://openfga.dev)**, offering simplified integration and internal conventions used within the Kopexa platform.

> Certain structural ideas and client abstractions are inspired by [Openlane’s](https://github.com/theopenlane/iam/blob/main/LICENSE) early approach to FGA integration.  
> All changes and extensions on top are independently developed and maintained by the Kopexa team.

### License

Licensed under **Apache 2.0**.  
See the [LICENSE](./LICENSE) file for details.

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
