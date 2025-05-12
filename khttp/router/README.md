# Router Package

This package provides a router with good defaults for Kopexa services. It is built on top of the `go-chi/chi` router and includes pre-configured middleware for observability, CORS, security headers, request ID, and metrics.

## Features

- Pre-configured router with good defaults
- Middleware registry for easy management of middleware
- Observability router for metrics and pprof endpoints
- CORS support
- Security headers
- Request ID middleware
- Metrics middleware

## Usage

### Basic Usage

```go
import "github.com/kopexa-grc/common/khttp/router"

func main() {
    r := router.New()
    // Add your routes here
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    // Start your server
}
```

### Observability Router

The observability router provides endpoints for metrics and pprof. It should not be exposed to the public internet.

```go
import "github.com/kopexa-grc/common/khttp/router"

func main() {
    r := router.NewObservabilityRouter()
    // Start your server
}
```

### Middleware Registry

You can customize the middleware registry to add or remove middleware.

```go
import "github.com/kopexa-grc/common/khttp/router"

func main() {
    customRegistry := router.NewMiddlewareRegistry()
    customRegistry.Register(cors.Middleware)
    r := router.NewWithOptions(router.WithMiddlewareRegistry(customRegistry))
    // Add your routes here
    // Start your server
}
```

## Configuration

The router can be configured using options. The following options are available:

- `WithMiddlewareRegistry`: Set a custom middleware registry

## License

This package is licensed under the BUSL-1.1 license. 