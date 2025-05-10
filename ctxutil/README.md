# ctxutil

The `ctxutil` package provides type-safe context value management using generics.
It allows storing and retrieving values from a context without type assertions
and with compile-time type safety.

## Features

- Type-safe storage and retrieval of context values
- Simple API for common use cases
- Support for default values and fallback functions
- Thread-safe through Go's built-in context system

## Installation

```bash
go get github.com/kopexa-grc/common/ctxutil
```

## Usage

### Storing Values in Context

```go
ctx := context.Background()
ctx = ctxutil.With(ctx, "my-value")
```

### Retrieving Values from Context

```go
// With error handling
if value, ok := ctxutil.From[string](ctx); ok {
    // Value found
}

// With panic on missing value
value := ctxutil.MustFrom[string](ctx)

// With default value
value := ctxutil.FromOr(ctx, "default-value")

// With fallback function
value := ctxutil.FromOrFunc(ctx, func() string {
    return "dynamic-default-value"
})
```

## Best Practices

1. **Type Safety**: Use the generic functions for type-safe access:
   ```go
   // Good
   value := ctxutil.From[string](ctx)
   
   // Bad
   value := ctx.Value("key").(string)
   ```

2. **Error Handling**: Use `From` with error handling for robust code:
   ```go
   if value, ok := ctxutil.From[string](ctx); ok {
       // Value found
   } else {
       // Value not found
   }
   ```

3. **Default Values**: Use `FromOr` or `FromOrFunc` for fallback values:
   ```go
   // Static default value
   value := ctxutil.FromOr(ctx, "default")
   
   // Dynamic default value
   value := ctxutil.FromOrFunc(ctx, func() string {
       return time.Now().Format(time.RFC3339)
   })
   ```

4. **Nested Contexts**: Be aware of context hierarchy:
   ```go
   ctx1 := ctxutil.With(ctx, "value1")
   ctx2 := ctxutil.With(ctx1, "value2")
   // ctx2 contains both values
   ```

## Performance

The package is optimized for high performance:
- Minimal allocations
- Efficient type checks
- Thread-safe operations

## License

BUSL-1.1 